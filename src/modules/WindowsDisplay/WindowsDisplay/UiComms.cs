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
	public class UiComms : BaseComms
	{
		public UiComms(MainWindow mw)
		{ 
			gui = mw; 
		}

		public void enqueue_query()
		{			
			lock (this)
			{
				Console.Out.WriteLine("Enqueing query");
				query_pending = true;
				Monitor.Pulse(this);
			}
		}

		public void enqueue_trigger( string tag )
		{
			lock (this)
			{
				trigger_pending = tag;
				Monitor.Pulse(this);
			}
		}

		public void enqueue_termination()
		{
			lock (this)
			{
				termination_required = true;
				Monitor.Pulse(this);
			}
		}
		/* Invoked as a thread from Program.cs */
		public void comms_loop()
		{
			Console.Out.Write("UiComms loop entered");
			bool qp;
			string tp;
			Console.Out.WriteLine("UI Comms entering event loop");
			while (true)
			{				
				Monitor.Enter(this);
				while (!query_pending && trigger_pending == null && !termination_required )
				{
					Monitor.Wait(this);
				}
				Console.Out.WriteLine("Something is pending");
				if (termination_required)
				{
					return;
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

		private void dispatch_query()
		{			
			SslStream peer = connect_to_server();
			Console.Out.WriteLine("Calling console to write query");
			Protocol.WriteQuery(peer);
			string[] tags = Protocol.ReadQueryResponse(peer);
			gui.DisplayQueryResponse(tags);
			peer.Dispose();
		}

		private void dispatch_trigger(string tag)
		{
			SslStream peer = connect_to_server();
			Protocol.WriteTrigger(peer, tag);
			Protocol.ReadTriggerResponse(peer);
			peer.Dispose();
		}

		SslStream connect_to_server()
		{
			
			Configuration conf = gui.GetCurrentConfig();
			if (conf == null)
			{
				Console.Error.WriteLine("No configuration defined, aborting connection");
				throw new ConnectionError();
			}
			TcpClient client = new TcpClient();
			Console.Out.WriteLine("Connecting to server " + conf.server.address + ":" + conf.server.port.ToString());
			client.Connect(conf.server.address, conf.server.port);
			SslStream peer = new SslStream(client.GetStream(), false);

			X509CertificateCollection certs = new X509CertificateCollection();
			Console.Out.WriteLine("Loading credentials from " + conf.server.cert_key + " pw '" + conf.server.password + "'");
			X509Certificate2 tmp = new X509Certificate2(conf.server.cert_key, conf.server.password);
			Console.Out.WriteLine("Loaded server cert, issuer " + tmp.Issuer + " subject " + tmp.Subject);
			Console.Out.WriteLine("Cert has key: " + tmp.HasPrivateKey);
			certs.Add(tmp);
			Console.Out.WriteLine("Pool contains " + certs.Count + " certificates");
			Console.Out.WriteLine("Certs loaded, starting authentication, hostname " + conf.server_hostname);
			peer.AuthenticateAsClient(conf.server_hostname,
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
			Console.Out.WriteLine("Connection to server established");
			return peer;
		}

		private bool query_pending = false;
		private string trigger_pending = null;
		private bool termination_required = false;
		private MainWindow gui = null;

	}

}
