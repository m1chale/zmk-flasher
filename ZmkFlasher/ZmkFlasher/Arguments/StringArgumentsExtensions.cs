using File = Dotcore.FileSystem.File;

namespace ZmkFlasher.Arguments;



internal static class StringArgumentsExtensions
{
    public static TypedArguments ToTypedArguments(this StringArguments args)
    => new (args.LeftFirmwarePath, args.RightFirmwarePath, bool.Parse(args.Verbose));
}
