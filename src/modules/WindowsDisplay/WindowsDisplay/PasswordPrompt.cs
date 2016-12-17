using System;
namespace WindowsDisplay
{
	public partial class PasswordPrompt : Gtk.Window
	{
		public PasswordPrompt(Configuration known) :
				base(Gtk.WindowType.Toplevel)
		{
			this.Build();
			valid = true;
			if (known.server.password != "")
			{
				server_pw_value.Text = known.server.password;
			}
			if (known.display.password != "")
			{
				display_pw_value.Text = known.display.password;
			}
		}

		public bool hasServerPassword()
		{
			return valid && server_pw_value.Text != "";
		}

		public string getServerPassword()
		{
			return server_pw_value.Text;
		}

		public bool hasDisplayPassword()
		{
			return valid && display_pw_value.Text != "";
		}

		public string getDisplayPassword()
		{
			return display_pw_value.Text;
		}

		protected void OnCancel(object sender, EventArgs e)
		{
			valid = false;
			this.Destroy();
		}

		protected void OnSave(object sender, EventArgs e)
		{
			valid = true;
			this.Destroy();
		}

		public bool valid = false;
	}
}
