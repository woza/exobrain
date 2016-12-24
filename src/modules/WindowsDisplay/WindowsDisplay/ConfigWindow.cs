using System;
using Gtk;
namespace WindowsDisplay
{
	public partial class ConfigWindow : Gtk.Window
	{
		public ConfigWindow() :
				base(Gtk.WindowType.Toplevel)
		{
			this.Build();
			new_conf = new Configuration(); /* Will load previous config if one exists */
			prefill();
		}

		private void prefill()
		{
			server_ck_value.Text = new_conf.server.cert_key;
			display_ck_value.Text = new_conf.display.cert_key;
			display_address_value.Text = new_conf.display.address;
			display_port_value.Text = new_conf.display.port.ToString();
			server_address_value.Text = new_conf.server.address;
			server_port_value.Text = new_conf.server.port.ToString();
			server_hostname_value.Text = new_conf.server_hostname;
			cmd_port_value.Text = new_conf.cmd_port.ToString();
		}

		public Configuration get_new_config()
		{
			return new_conf;
		}

		protected void CancelReconfig(object sender, EventArgs e)
		{
			new_conf = null;
			this.Destroy();
		}

		protected void ApplyReconfig(object sender, EventArgs e)
		{
			new_conf.server.cert_key = server_ck_value.Text;
			new_conf.display.cert_key = display_ck_value.Text;

			new_conf.display.address = display_address_value.Text;
			new_conf.display.port = Convert.ToInt16(display_port_value.Text);
			new_conf.server.address = server_address_value.Text;
			new_conf.server.port = Convert.ToInt16(server_port_value.Text);

			new_conf.server_hostname = server_hostname_value.Text;
			new_conf.cmd_port = Convert.ToInt32(cmd_port_value.Text);
			new_conf.Save(); /* Save before validation so that a validation failure doesn't imply re-entry from scratch */
			if (!new_conf.validate())
			{
				MessageDialog oops = new MessageDialog(this,
													   DialogFlags.DestroyWithParent | DialogFlags.Modal,
													   MessageType.Error,
													   ButtonsType.Close, "Invalid configuration: " + new_conf.explain());
				oops.Run();
				oops.Destroy();
			}
			this.Destroy();
		}

		protected void OnBrowseServerCred(object sender, EventArgs e)
		{
			get_file(server_ck_value);
		}

		protected void OnBrowseDisplayCred(object sender, EventArgs e)
		{
			get_file(display_ck_value);
		}

		private void get_file(Entry dest)
		{
			FileChooserDialog question = new FileChooserDialog("Select file...",
															  this,
															   FileChooserAction.Open,
															   "Cancel", ResponseType.Cancel,
															   "Select", ResponseType.Accept);
			int outcome = question.Run();
			if (outcome == (int)ResponseType.Accept)
			{
				dest.Text = question.Filename;
			}
			question.Destroy();
		}

		private Configuration new_conf = null;

		protected void OnBrowseDisplayCert(object sender, EventArgs e)
		{
		}
	}
}
