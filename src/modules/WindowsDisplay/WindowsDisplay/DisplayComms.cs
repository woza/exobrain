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
using Gtk;

namespace WindowsDisplay
{
	public class DisplayComms 
	{
		public DisplayComms(MainWindow mw)
		{ 
			gui = mw;
		}

		public void enqueue_termination()
		{
			try
			{
				NetworkStream ns = connect_to_back_channel();
				Protocol.WriteQuit(ns);
				ns.Dispose();
			}
			catch (SocketException e)
			{
				Console.Error.WriteLine("Error encountered while trying to terminate: " + e.ToString());
			}
		}

		private NetworkStream connect_to_back_channel()
		{
			Configuration conf = gui.GetCurrentConfig();
			IPEndPoint dest;
			lock (this)
			{
				dest = new IPEndPoint(IPAddress.Loopback, conf.cmd_port);
			}
			System.Net.Sockets.Socket sock = new System.Net.Sockets.Socket(SocketType.Stream, ProtocolType.Tcp);
			sock.Connect(dest);
			return new NetworkStream(sock, true);
		}

		public bool is_sufficient(Configuration conf)
		{
			lock(this)
			{
				if (conf == null)
				{
					return false;
				}
				if (conf.cmd_port == 0)
				{
					return false;
				}
				if (conf.display == null)
				{
					return false;
				}
				if (conf.display.address == "")
				{
					return false;
				}
				if (conf.display.port == 0)
				{
					return false;
				}

				if (conf.display.cert_key == "")
				{
					return false;
				}
			}
			Console.Error.WriteLine("Haveconfig passed");
			return true;
		}

		/* Invoked as a thread from Program.cs */
		public void comms_loop()
		{
			while (true)
			{
				Stream client = next_client();
				Tuple<UInt32, Protocol.cmd_t> header = Protocol.get_display_message_header(client);
				switch (header.Item2)
				{
					case Protocol.cmd_t.CMD_QUIT:
						return;
					case Protocol.cmd_t.CMD_DISPLAY:
						string pw = Protocol.ReadDisplayedTag(client, header.Item1);
						gui.DisplayPassword(pw);
						Protocol.WriteDisplayResponse(client);
						break;
				}
				client.Dispose();
			}					
		}

		private Stream next_client()
		{
			while (true)
			{
				Configuration curr_config = gui.GetCurrentConfig();
				if (!curr_config.Equals(last_config))
				{
					if (remote_server != null)
					{
						remote_server.Dispose();
					}
					remote_server = new System.Net.Sockets.Socket(SocketType.Stream, ProtocolType.Tcp);
					remote_server.Bind(new IPEndPoint(IPAddress.Any, curr_config.display.port));
					remote_server.Listen(5);
					if (back_channel != null)
					{
						back_channel.Dispose();
					}
					back_channel = new System.Net.Sockets.Socket(SocketType.Stream, ProtocolType.Tcp);
					back_channel.Bind(new IPEndPoint(IPAddress.Loopback, curr_config.cmd_port));
					back_channel.Listen(5);
				}
				last_config = curr_config;
				ArrayList empty = new ArrayList();
				ArrayList acceptors = new ArrayList();

				acceptors.Add(remote_server);
				acceptors.Add(back_channel);
				System.Net.Sockets.Socket.Select(acceptors, empty, empty, -1);
				/* Back channel takes priority */
				if (acceptors.Contains(back_channel))
				{
					System.Net.Sockets.Socket peer = back_channel.Accept();
					return new NetworkStream(peer, true);
				}
				if (acceptors.Contains(remote_server))
				{
					System.Net.Sockets.Socket peer = remote_server.Accept();
					try
					{
						return upgrade_to_tls(peer);
					}
					catch (IOException)
					{
						/* Problem with this client, ignore and try again */
						continue;
					}
					catch(AuthenticationException)
					{
						/* Problem with this client, ignore and try again */
						continue;
					}

				}
				Console.Error.WriteLine("Select returned with no readable sockets");
				throw new ConnectionError();
			}
		}

		private SslStream upgrade_to_tls(System.Net.Sockets.Socket raw)
		{
			SslStream ret = new SslStream(new NetworkStream(raw, true), false);
			X509Certificate2 cert = new X509Certificate2(last_config.display.cert_key,
													   last_config.display.password);
			ret.AuthenticateAsServer(cert, true, SslProtocols.Tls12, false);
			if (!ret.IsSigned)
			{
				Console.Error.WriteLine("Failed to establish signed server-side socket");
				throw new ConnectionError();
			}
			if (!ret.IsEncrypted)
			{
				Console.Error.WriteLine("Failed to establish encrypted server-side socket");
				throw new ConnectionError();
			}
			if (!ret.IsMutuallyAuthenticated)
			{
				Console.Error.WriteLine("Failed to establish mutually authenticated server-side socket");
				throw new ConnectionError();
			}

			return ret;
		}

		Configuration last_config = null;
		private MainWindow gui = null;
		private System.Net.Sockets.Socket remote_server = null;
		private System.Net.Sockets.Socket back_channel = null;
	}
}


