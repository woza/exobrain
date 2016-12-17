using System;
using System.IO;

namespace WindowsDisplay
{
	public class Configuration : IEquatable<Configuration>
	{
		public Configuration()
		{
			server = new ConnInfo("server");
			display = new ConnInfo("display");
		
			/* Try and populate from previous settings */
			string dir = Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData);
			string full_path = Path.Combine(dir, "exobrain.conf");
			if (File.Exists(full_path))
			{
				Console.Out.WriteLine("Reading configuration from " + full_path);
				string[] old = File.ReadAllLines(full_path);
				fromStrings(old);
			}
		}

		public bool Equals(Configuration other)
		{
			if (other == null)
			{
				return false;
			}
			return server.Equals(other.server) &&
						 display.Equals(other.display) &&
						 cmd_port == other.cmd_port &&
						 server_hostname.Equals(other.server_hostname);
		}

		public void Save()
		{
			string dir = Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData);
			string full_path = Path.Combine(dir, "exobrain.conf");
			File.WriteAllText(full_path, toString());
		}

		public bool validate()
		{
			Console.Out.WriteLine("=====\nValidating " + this.toString()+"=====");
			return true;
		}

		public string explain()
		{
			return validation_failure_reason;
		}

		public class ConnInfo : IEquatable<ConnInfo>
		{
			string tag;
			public string address;
			public int port;
			public string cert_key;
			public string ca;
			public string password;
			public ConnInfo(string t)
			{
				tag = t;
				address = "";
				port = 0;
				cert_key = "";
				ca = "";
				password = "";
			}

			public bool Equals(ConnInfo other)
			{
				if (other == null)
				{
					return false;
				}
				return tag.Equals(other.tag) &&
						  address.Equals(other.address) &&
									   port == other.port &&
						  cert_key.Equals(other.cert_key) &&
						  ca.Equals(other.ca) &&
						  password.Equals(other.password);
			}


			public string toString()
			{
				string ret = "";
				ret += tag + ".address = " + address + "\n";
				ret += tag + ".port = " + port.ToString() + "\n";
				ret += tag + ".cert_key = " + cert_key + "\n";
				ret += tag + ".ca = " + ca + "\n";
				return ret;
			}
		};

		/* Don't store passwords to disk */
		public string toString()
		{
			string ret = "";

			ret += server.toString();
			ret += display.toString();
			ret += "server_hostname = " + server_hostname + "\n";
			ret += "Validation failure = " + validation_failure_reason + "\n";
			ret += "cmd_port = " + cmd_port + "\n";
			return ret;
		}

		private void fromStrings(string[] src)
		{
			foreach (string line in src)
			{
				string[] bits = line.Split('=');
				string key = bits[0].Trim();
				string value = bits[1].Trim();
				if (key == "server.cert_key")
				{
					server.cert_key = value;
					continue;
				}
				if (key == "server.ca")
				{
					server.ca = value;
					continue;
				}
				if (key == "server.address")
				{
					server.address = value;
					continue;
				}
				if (key == "server.port")
				{
					server.port = Convert.ToUInt16(value);
					continue;
				}
				if (key == "display.cert_key")
				{
					display.cert_key = value;
					continue;
				}
				if (key == "display.ca")
				{
					display.ca = value;
					continue;
				}
				if (key == "display.address")
				{
					display.address = value;
					continue;
				}
				if (key == "display.port")
				{
					display.port = Convert.ToUInt16(value);
					continue;
				}

				if (key == "server_hostname")
				{
					server_hostname = value;
					continue;
				}
				if (key == "cmd_port")
				{
					cmd_port = Convert.ToInt32(value);
					continue;
				}
			}
		}

		public ConnInfo server;
		public ConnInfo display;
		private string validation_failure_reason = "";
		public string server_hostname { get; set; }
		public int cmd_port = 0;
	}
}
