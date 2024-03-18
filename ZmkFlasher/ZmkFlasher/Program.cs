
using Dotcore.FileSystem.File;
using ZmkFlasher.Arguments;
using ZmkFlasher.WaitRemovableDevice;

var result = CommandLine.Parser.Default.ParseArguments<StringArguments>(args);
if (result.Errors.Any()) throw new Exception($"invalid arguments {result}");
var arguments = result.Value.ToTypedArguments();


await WaitAndCopy.WaitAndCopyFirmware(("GLV80LHBOOT", arguments.LeftFirmware), ("GLV80RHBOOT", arguments.RightFirmware));
Console.WriteLine("Done");