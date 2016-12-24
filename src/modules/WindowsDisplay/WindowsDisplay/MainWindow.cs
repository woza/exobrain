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
		Log("Configuration loaded from " + curr_config.src_path);
		ui_comms_thread = new Thread(new ThreadStart(ui_comms_worker.comms_loop));
		ui_comms_thread.Start();

		display_comms_thread = new Thread(new ThreadStart(display_comms_worker.comms_loop));
		if (display_comms_worker.is_sufficient(curr_config))
		{
			display_comms_thread.Start();
		}

		setup_tag_display();

	}

	private void setup_tag_display()
	{
		view = new ComboBox();

		view.Changed += Trigger;
		main_grid.Attach(view, 0, 1, 2, 3);

		tag_list = new ListStore(typeof(string));
		view.Model = tag_list;

		CellRenderer render = new CellRendererText();
		view.PackStart(render, true);
		view.AddAttribute(render, "text", 0);
		ShowAll();
	}

	public void Log(string msg)
	{
		log_output.Buffer.Text += msg + "\n";
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
		Array.Sort(tags);
		tag_list.Clear();
		string[] element = new string[1];
		for (uint i = 0; i < tags.Length; ++i)
		{
			element[0] = tags[i];
			tag_list.AppendValues(element);
		}
		view.ShowAll();
		Log("New tag set received, make selection from drop-down box");
	}

	protected void Trigger(object sender, EventArgs a)
	{
		ComboBox src = (ComboBox)sender;
		string tag = src.ActiveText;
		ui_comms_worker.enqueue_trigger(tag);
	}

	protected void OnDeleteEvent(object sender, DeleteEventArgs a)
	{
		Application.Quit();
		a.RetVal = true;
	}

	protected void DisplayReconfigure(object sender, EventArgs e)
	{
		config_window = new ConfigWindow();
		config_window.Destroyed += OnReconfigure;
		config_window.Show();
	}

	protected void OnReconfigure(object sender, EventArgs e)
	{
		Configuration pending_config = config_window.get_new_config();
		if (pending_config != null)
		{
			lock(this)
			{
				curr_config = pending_config;
			}
		}

		if ( display_comms_thread.ThreadState == ThreadState.Unstarted && display_comms_worker.is_sufficient(pending_config))
		{
			display_comms_thread.Start();
		}
	}

	protected void OnRefresh(object sender, EventArgs e)
	{
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
		ui_comms_worker.enqueue_termination();
		display_comms_worker.enqueue_termination();
	}
	private ListStore tag_list;

	private Configuration curr_config = null;
	private ConfigWindow config_window;
	private Thread display_comms_thread = null;

	private DisplayComms display_comms_worker = null;
	private Thread ui_comms_thread = null;
	private UiComms ui_comms_worker = null;
	private ComboBox view = null;

	protected void OnTagSelected(object sender, EventArgs e)
	{
	}
}
