using System;
namespace WindowsDisplay
{
	public partial class PasswordPrompt : Gtk.Window
	{
		public PasswordPrompt() :
				base(Gtk.WindowType.Toplevel)
		{
			this.Build();
		}
	}
}
