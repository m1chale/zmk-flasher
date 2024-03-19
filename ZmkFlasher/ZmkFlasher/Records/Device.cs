using Dotcore.FileSystem.Directory;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ZmkFlasher.Records;

public record Device(string Name, string Label, Info[] MountPoints);
