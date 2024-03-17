
using CommandLine;
using File = Dotcore.FileSystem.File;


internal record TypedArguments(File.Info LeftFirmwarePath, File.Info RightFirmwarePath);
