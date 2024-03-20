
using CommandLine;

internal class StringArguments
{
    [Option('l', "left", Required = true, HelpText = "Path to the left firmware file")]
    public string LeftFirmwarePath { get; set; } = string.Empty;
    [Option('r', "right", Required = true, HelpText = "Path to the right firmware file")]
    public string RightFirmwarePath { get; set; } = string.Empty;

    [Option('d', "dryRun", Required = false, Default = "false", HelpText = "DryRun")]
    public string DryRun { get; set; } = bool.FalseString;
}

