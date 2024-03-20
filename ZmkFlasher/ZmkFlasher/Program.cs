
using Dotcore.FileSystem.File;
using ZmkFlasher;
using ZmkFlasher.Arguments;
using ZmkFlasher.BL;
using ZmkFlasher.Lib;
using ZmkFlasher.WaitRemovableDevice;

var result = CommandLine.Parser.Default.ParseArguments<StringArguments>(args);
if (result.Errors.Any()) throw new Exception($"invalid arguments {result}");
var arguments = result.Value.ToTypedArguments();
DryRunExtensions.IsDryRun = arguments.DryRun;

arguments.LeftFirmware.ThrowIfNotExistsOrDryRun();
arguments.RightFirmware.ThrowIfNotExistsOrDryRun();

Console.WriteLine("Connect left and right bootloader");
var leftDeviceTask = IWaitForDevice.Instance.WaitForDevice("GLV80LHBOOT");
var rightDeviceTask = IWaitForDevice.Instance.WaitForDevice("GLV80RHBOOT");
var leftAndRightDevices = await Task.WhenAll(leftDeviceTask, rightDeviceTask);
if(leftAndRightDevices.Length != 2) throw new Exception("Failed to find left and right bootloader");

var leftDevice = leftAndRightDevices[0];
var rightDevice = leftAndRightDevices[1];
var leftMountPoint = leftDevice.MountPoints.Single();
var rightMountPoint = rightDevice.MountPoints.Single();

Console.WriteLine("copying firmware to left bootloader");
arguments.LeftFirmware.CopyToOrDryRun(leftMountPoint);
Console.WriteLine("copying firmware to right bootloader");
arguments.RightFirmware.CopyToOrDryRun(rightMountPoint);

Console.WriteLine("cleaning up");
await ICleanUpDevice.Instance.CleanUp(leftDevice);
await ICleanUpDevice.Instance.CleanUp(rightDevice);
Console.WriteLine("done");