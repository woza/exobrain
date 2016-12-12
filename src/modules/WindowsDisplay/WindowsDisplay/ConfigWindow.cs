using System;
namespace WindowsDisplay
{
	public partial class ConfigWindow : Gtk.Window
	{
		public ConfigWindow() :
				base(Gtk.WindowType.Toplevel)
		{
			this.Build();
		}
	}
}
