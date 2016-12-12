using System;
using System.Collections;
using System.Net;
using System.Net.Sockets;
using System.Net.Security;
using System.Security.Authentication;
using System.Text;
using System.Security.Cryptography.X509Certificates;
using System.IO;
using System.Threading;

namespace WindowsDisplay
{
	public class DisplayComms
	{
		public DisplayComms()
		{
			cert = X509Certificate.CreateFromCertFile(cert_path);
		}

		public void enqueue_new_config(Configuration conf)
		{
			lock (this)
			{
				pending_config = conf;
			}
		}

		public void enqueue_query()
		{
			lock (this)
			{
				query_pending = true;
			}
			Monitor.Pulse(this);
		}

		public void enqueue_trigger( string tag )
		{
			lock (this)
			{
				trigger_pending = tag;
			}
			Monitor.Pulse(this);
		}

		public void comms_loop()
		{
			bool qp;
			string tp;
			while (true)
			{				
				
				ensure_config_current();
				Monitor.Enter(this);
				while (!query_pending && trigger_pending == null)
				{
					Monitor.Wait(this);
				}
				qp = query_pending;
				tp = trigger_pending;
				query_pending = false;
				trigger_pending = null;
				Monitor.Exit(this);
				if (qp)
				{
					Console.Out.WriteLine("Dispatching query");
					dispatch_query();
				}
				if (tp != null)
				{
					Console.Out.WriteLine("Dispatching trigger " + trigger_pending);
					dispatch_trigger(tp);
				}
			}
		}

		private void ensure_config_current()
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
			if (need_swap)
			{
				if (server != null)
				{
					server.Stop();
				}
				server = new TcpListener(IPAddress.Any, Int32.Parse(active_config.port));
				server.Start();
			}
		}

		private void dispatch_query()
		{
			SslStream peer = connect_to_server();
			Protocol.WriteQuery(peer);
			string[] tags = Protocol.ReadQueryResponse(peer);
			gui.DisplayQueryResponse(tags);
			peer.Dispose();
		}

		private void dispatch_trigger(string tag)
		{
			SslStream peer = connect_to_server();
			Protocol.WriteTrigger(peer, tag);
			string pw = Protocol.ReadTriggerResponse(peer);
			gui.DisplayTriggerResponse(pw);
			peer.Dispose();
		}

		SslStream connect_to_server()
		{
			TcpClient client = new TcpClient();
			client.Connect(active_config.server, active_config.port);
			SslStream peer = new SslStream(client.GetStream(), false);

			X509CertificateCollection certs = new X509CertificateCollection();
			certs.Add(new X509Certificate(active_config.tls_combined_path, active_config.key_password));
			peer.AuthenticateAsClient(active_config.server_hostname,
									  certs,
									  SslProtocols.Tls12,
									  false);
			if (!peer.IsEncrypted)
			{
				Console.Error.WriteLine("Server connection not encrypted, aborting");
				peer.Dispose();
				throw new ConnectionError();
			}
			if ( !peer.IsSigned)
			{
				Console.Error.WriteLine("Server connection not signed, aborting");
				peer.Dispose();
				throw new ConnectionError();
			}
			if ( !peer.IsMutuallyAuthenticated )
			{
				Console.Error.WriteLine("Server connection not mutually authenticated, aborting");
				peer.Dispose();
				throw new ConnectionError();
			}
			return peer;
		}

		private bool query_pending = false;
		private string trigger_pending = null;

		
		private Configuration active_config;
		private Configuration pending_config = null;
		private TcpListener server = null;

		private X509Certificate cert = null;

		private MainWindow gui = null;

		private class ConnectionError : Exception{ }
	}

}
