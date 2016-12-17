using System;
using System.Net.Security;
namespace WindowsDisplay
{
	public class BaseComms
	{
		public BaseComms()
		{
		}

		public void enqueue_new_config(Configuration conf)
		{
			lock (this)
			{
				pending_config = conf;
			}
		}

		protected bool ensure_config_current()
		{
			bool need_swap = false;
			lock (this)
			{
				if (pending_config != null)
				{
					need_swap = true;
					active_config = pending_config;
					pending_config = null;
				}
			}
			return need_swap;
		}

		protected UInt32 get_u32(SslStream src)
		{
			byte[] raw = new byte[4];
			src.Read(raw, 0, 4);
			if (BitConverter.IsLittleEndian)
			{
				Array.Reverse(raw);
			}
			return BitConverter.ToUInt32(raw, 0);
		}

		protected Configuration active_config;
		private Configuration pending_config = null;

		protected class ConnectionError : Exception { }
	}
}
