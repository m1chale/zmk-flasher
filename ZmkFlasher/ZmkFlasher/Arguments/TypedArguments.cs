
using CommandLine;
using File = Dotcore.FileSystem.File;


internal record TypedArguments(File.Info LeftFirmware, File.Info RightFirmware);
