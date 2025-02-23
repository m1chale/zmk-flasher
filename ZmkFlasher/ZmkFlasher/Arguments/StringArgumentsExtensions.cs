using File = Dotcore.FileSystem.File;

namespace ZmkFlasher.Arguments;



internal static class StringArgumentsExtensions
{
    public static TypedArguments ToTypedArguments(this StringArguments args)
    {
        if (!string.IsNullOrWhiteSpace(args.LeftFirmwarePath) && !string.IsNullOrWhiteSpace(args.LeftAndRightFirmwarePath)) throw new Exception("LeftFirmwarePath and LeftAndRightFirmwarePath are mutually exclusive");
        if (!string.IsNullOrWhiteSpace(args.RightFirmwarePath) && !string.IsNullOrWhiteSpace(args.LeftAndRightFirmwarePath)) throw new Exception("RightFirmwarePath and LeftAndRightFirmwarePath are mutually exclusive");

        Console.WriteLine($"LeftFirmwarePath: {args.LeftFirmwarePath}");
        Console.WriteLine($"RightFirmwarePath: {args.RightFirmwarePath}");
        Console.WriteLine($"LeftAndRightFirmwarePath: {args.LeftAndRightFirmwarePath}");

        var leftFirmwarePath = string.IsNullOrWhiteSpace(args.LeftFirmwarePath) ? args.LeftAndRightFirmwarePath : args.LeftFirmwarePath;
        var rightFirmwarePath = string.IsNullOrWhiteSpace(args.RightFirmwarePath) ? args.LeftAndRightFirmwarePath : args.RightFirmwarePath;

        if (string.IsNullOrWhiteSpace(leftFirmwarePath)) throw new Exception("LeftFirmwarePath is required");
        if (string.IsNullOrWhiteSpace(rightFirmwarePath)) throw new Exception("RightFirmwarePath is required");

        return new TypedArguments(new File.Info(leftFirmwarePath), new File.Info(rightFirmwarePath), bool.Parse(args.DryRun));
    }
}
