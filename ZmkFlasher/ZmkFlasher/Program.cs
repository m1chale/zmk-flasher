
using Dotcore.FileSystem.File;
using ZmkFlasher.Arguments;
using ZmkFlasher.WaitRemovableDevice;

var result = CommandLine.Parser.Default.ParseArguments<StringArguments>(args);
if (result.Errors.Any()) throw new Exception($"invalid arguments {result}");
var arguments = result.Value.ToTypedArguments();

Console.WriteLine("Connect left keyboard");
var directory = await Wait.ForDevice("GLV80LHBOOT");
arguments.LeftFirmwarePath.CopyTo(directory);

Console.WriteLine("Connect right keyboard");
directory = await Wait.ForDevice("GLV80RHBOOT");
arguments.RightFirmwarePath.CopyTo(directory);

