using System.Net.Security;
using System;
using System.Text;
using System.IO;
namespace WindowsDisplay
{
	public class Protocol
	{
		public Protocol()
		{
		}

		private static void WriteParameterlessMessage(Stream dest, uint msg)
		{
			const uint msg_size = 4;
			put_u32(dest, msg_size);
			put_u32(dest, msg);
		}

		public static void WriteQuit(Stream dest)
		{
			WriteParameterlessMessage(dest, (uint)cmd_t.CMD_QUIT);
		}

		public static void WriteReconfig(Stream dest)
		{
			WriteParameterlessMessage(dest, (uint)cmd_t.CMD_RECONF);
		}

		public static void WriteQuery(Stream dest)
		{
			WriteParameterlessMessage(dest, (uint)cmd_t.CMD_QUERY_ALL);
		}

		public static void WriteFailure(Stream dest)
		{
			WriteParameterlessMessage(dest, (uint)status_t.STATUS_FAIL);
		}

		public static void WriteDisplayResponse(Stream dest)
		{
			WriteParameterlessMessage(dest, (uint)status_t.STATUS_OK);
		}

		public static string[] ReadQueryResponse(Stream src)
		{
			UInt32 msg_size = get_u32(src);
			status_t status = (status_t)get_u32(src);
			if (status != status_t.STATUS_OK)
			{
				throw new OperationFailed();
			}
			Console.Out.WriteLine("Query response STATUS_OK");
			encode_t coding = (encode_t)get_u32(src);
			uint count = get_u32(src);
			Console.Out.WriteLine("Count = " + count);
			string[] ret = new string[count];
			for (uint i = 0; i < count; ++i)
			{
				uint sz = get_u32(src);
				byte[] raw_string = new byte[sz];
				Int32 got = src.Read(raw_string, 0, (int)sz);
				if (got != sz)
				{
					Console.Error.WriteLine("Tried to read " + sz + " bytes but only managed " + got);
					throw new OperationFailed();
				}
				ret[i] = decode_string(raw_string, coding);
				Console.Out.WriteLine("Query response [" + i + "] sz " + sz + " = " + ret[i]);
			}
			return ret;
		}

		public static void WriteTrigger(Stream dest, string tag)
		{

			const cmd_t cmd = cmd_t.CMD_TRIGGER;
			const encode_t coding = encode_t.ENCODE_UTF8;
			byte[] coded_tag = Encoding.UTF8.GetBytes(tag);
			uint msg_size = (uint)(4 + 4) + (uint)coded_tag.Length;

			put_u32(dest, msg_size);
			put_u32(dest, (uint)cmd);
			put_u32(dest, (uint)coding);
			dest.Write(coded_tag, 0, coded_tag.Length);
		}

		public static void ReadTriggerResponse(Stream src)
		{
			status_t status = (status_t)get_u32(src);
			if (status != status_t.STATUS_OK)
			{
				throw new OperationFailed();
			}
		}

		public static Tuple<UInt32, cmd_t> get_display_message_header(Stream src)
		{
			UInt32 sz = get_u32(src);
			Console.Out.WriteLine("Display message header size: " + sz);
			cmd_t cmd = (cmd_t)get_u32(src);
			Console.Out.WriteLine("Command = " + cmd);
			return Tuple.Create(sz, cmd);
		}

		public static string ReadDisplayedTag(Stream src, UInt32 payload_size)
		{
			/* Command and size have already been extracted from the stream */
			payload_size -= 8;
			encode_t coding = (encode_t)get_u32(src);
			payload_size -= 4;
			byte[] raw = new byte[payload_size];
			src.Read(raw, 0, (int)payload_size);
			return decode_string(raw, coding);
		}

		private static string decode_string(byte[] raw, encode_t src_coding)
		{
			switch (src_coding)
			{
				case encode_t.ENCODE_UTF8:
					Console.Out.WriteLine("Decoding UTF8 string");
					return Encoding.UTF8.GetString(raw);
				case encode_t.ENCODE_ASCII:
					Console.Out.WriteLine("Decoding ASCII string");
					return Encoding.ASCII.GetString(raw);
				default:
					throw new InvalidEncoding();
			}
		}

		private static void put_u32(Stream dest, UInt32 val )
		{			
			byte[] buffer = BitConverter.GetBytes(val);
			Console.Out.WriteLine("Input value " + val.ToString() + " bytes " + BitConverter.ToString(buffer));
			if (BitConverter.IsLittleEndian)
			{
				Console.Out.WriteLine("Switching byte order");
				/* Convert from network-byte-order */
				Array.Reverse(buffer);
			}
			Console.Out.WriteLine("Sending bytes " + BitConverter.ToString(buffer));
			dest.Write(buffer, 0, buffer.Length);
		}

		private static UInt32 get_u32(Stream src)
		{
			byte[] buffer = new byte[4];
			src.Read(buffer, 0, 4);
			if (BitConverter.IsLittleEndian)
			{
				/* Convert from network-byte-order */
				Array.Reverse(buffer);
			}
			return BitConverter.ToUInt32(buffer, 0);
		}

		private enum status_t : uint { STATUS_OK = 0, STATUS_FAIL };
		public enum cmd_t : uint { CMD_QUERY_ALL = 0, CMD_TRIGGER, CMD_DISPLAY, CMD_RECONF=500, CMD_QUIT=501 };
		private enum encode_t : uint { ENCODE_ASCII = 0, ENCODE_UTF8 };
			
	}
	public class OperationFailed : Exception { }
    public class InvalidEncoding : Exception { }
}
