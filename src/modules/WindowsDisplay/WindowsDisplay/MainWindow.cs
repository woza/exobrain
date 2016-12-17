using System;
using Gtk;
using WindowsDisplay;
using System.Threading;

public partial class MainWindow : Gtk.Window
{
	public MainWindow() : base(Gtk.WindowType.Toplevel)
	{
		Build();
		this.DeleteEvent += GracefulShutdown;
		display_comms_worker = new DisplayComms(this);
		ui_comms_worker = new UiComms(this);

		/* Initial config from existing settings */
		curr_config = new Configuration();


		ui_comms_worker.enqueue_new_config(curr_config);
		ui_comms_thread = new Thread(new ThreadStart(ui_comms_worker.comms_loop));
		ui_comms_thread.Start();

		display_comms_worker.enqueue_new_config(curr_config);
		display_comms_thread = new Thread(new ThreadStart(display_comms_worker.comms_loop));
		if (display_comms_worker.have_config())
		{
			display_comms_thread.Start();
		}
	}

	public void Log(string msg)
	{
		Console.Out.WriteLine(msg);
	}

	public void DisplayQueryResponse(string[] tags)
	{
		/* Run this on the main GUI thread */
		Gtk.Application.Invoke(delegate { draw_tag_list(tags); });
		Log("Query response received");
	}

	public void DisplayPassword(string pw)
	{
		/* Run this on the main GUI thread */
		Gtk.Application.Invoke(delegate { pw_label.Text = pw; } );
		Log("Password retrieved");
	}

	public Configuration GetCurrentConfig()
	{
		Configuration ret;
		lock (this)
		{
			ret = curr_config;
		}
		return ret;
	}
			
	private void draw_tag_list(string[] tags)
	{
		Console.Out.WriteLine("Drawing tag list (length " + tags.Length + ") onto screen");
		if (tag_display != null)
		{
			foreach (Widget w in tag_display.AllChildren)
			{
				w.Destroy();
			}
			tag_display.Destroy();
		}
			
		tag_display = new Table((uint)tags.Length, 1, false);
		for (uint i = 0; i < tags.Length; ++i)
		{
			Console.Out.WriteLine("Drawing tag [" + i + "] = " + tags[i]);
			Button b = new Button(tags[i]);
			b.Clicked += Trigger;
			tag_display.Attach(b, 0, 1, i, i + 1);
			b.Show();
		}
		main_grid.ShowAll();
		main_grid.Attach(tag_display, 0, 1, 2, 3);
		tag_display.Show();
	}

	protected void Trigger(object sender, EventArgs a)
	{
		Button src = (Button)sender;
		string tag = src.Label;
		Console.Out.WriteLine("GUI triggering tag " + tag);
		ui_comms_worker.enqueue_trigger(tag);
	}

	protected void OnDeleteEvent(object sender, DeleteEventArgs a)
	{
		Application.Quit();
		a.RetVal = true;
	}

	protected void DisplayReconfigure(object sender, EventArgs e)
	{
		Console.Out.WriteLine("Displaying config information");
		config_window = new ConfigWindow();
		config_window.Destroyed += OnReconfigure;
		config_window.Show();
	}

	protected void OnReconfigure(object sender, EventArgs e)
	{
		Console.Out.WriteLine("Handling reconfigure input");
		Configuration pending_config = config_window.get_new_config();
		if (pending_config != null)
		{
			lock(this)
			{
				curr_config = pending_config;
			}
		}
		display_comms_worker.enqueue_new_config(curr_config);
		if ( display_comms_thread.ThreadState == ThreadState.Unstarted && display_comms_worker.have_config())
		{
			display_comms_thread.Start();
		}
	}

	protected void OnRefresh(object sender, EventArgs e)
	{
		Console.Out.WriteLine("Refresh click detected");
		lock(this)
		{
			if (curr_config.server.password == "" || curr_config.display.password == "" )
			{				
				PasswordPrompt pp = new PasswordPrompt(curr_config);
				pp.Destroyed += UpdatePassword;
				pp.Show();
				return;
			}
		}
		ActualRefresh();
	}

	protected void UpdatePassword(object sender, EventArgs e)
	{
		PasswordPrompt pp = (PasswordPrompt)sender;
		lock(this)
		{
			if (pp.hasServerPassword())
			{
				curr_config.server.password = pp.getServerPassword();
			}
			if (pp.hasDisplayPassword())
			{
				curr_config.display.password = pp.getDisplayPassword();
			}
		}
		pp.Dispose();
		ActualRefresh();
	}


	protected void ActualRefresh()
	{
		ui_comms_worker.enqueue_query();
	}

	protected void OnClear(object sender, EventArgs e)
	{
		pw_label.Text = "";
	}

	protected void GracefulShutdown(object sender, EventArgs e)
	{
		Console.Out.WriteLine("Graceful shutdown invoked");
		ui_comms_worker.enqueue_termination();
		display_comms_worker.enqueue_termination();
	}
	private Configuration curr_config = null;
	private ConfigWindow config_window;
	private Thread display_comms_thread = null;

	private DisplayComms display_comms_worker = null;
	private Thread ui_comms_thread = null;
	private UiComms ui_comms_worker = null;
	private Table tag_display = null;
}
