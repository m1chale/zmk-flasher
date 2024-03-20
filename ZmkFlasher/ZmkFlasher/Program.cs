
using Dotcore.FileSystem.File;
using ZmkFlasher.Arguments;
using ZmkFlasher.Lib;
using ZmkFlasher.WaitRemovableDevice;

var result = CommandLine.Parser.Default.ParseArguments<StringArguments>(args);
if (result.Errors.Any()) throw new Exception($"invalid arguments {result}");
var arguments = result.Value.ToTypedArguments();
arguments.LeftFirmware.ThrowIfNotExists();
arguments.RightFirmware.ThrowIfNotExists();

Console.WriteLine("Connect left bootloader");
await IWaitAndCopy.Instance.WaitForDeviceAndCopy("GLV80LHBOOT", arguments.LeftFirmware);
Console.WriteLine("Connect right bootloader");
await IWaitAndCopy.Instance.WaitForDeviceAndCopy("GLV80RHBOOT", arguments.RightFirmware);
Console.WriteLine("Done");