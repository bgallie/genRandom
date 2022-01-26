# Name

**genRandom** - generate psudo-random data using the tnt2engine random number generator.

## Description

**genRandom** generates a stream of psudo-random data using the *tnt2engine* psudo-randmon number generator.  The data can be output as:

- a stream of ASCII '0' and '1' characters.
- a stream of bytes.
- a stream of hexadecimal encoded bytes.
- a series of tab seperated XY coordinates in a rectangle of a given size.

## Usage:

        genRandom [command]

### Available Commands:

      binary      Generate a stream of (psudo)random ASCII '0' and '1' characters.
      completion  Generate the autocompletion script for the specified shell
      data        Generate a stream of (psudo)random bytes.
      help        Help about any command
      hex         Generate a stream of (psudo)random hexadecimal encoded bytes.
      points      Generate a series of (psudo)random X,Y coordinates in a rectanglee of a given size
      version     Display version information

### Flags:

          --config string         config file (default is $HOME/.tnt2.yaml)
          --count=N string        starting block count.
      -h, --help                  help for genRandom
      -o, --outputFile string     Name of the file containing the generated (psudo)random data. (default "-")
      -f, --proformafile string   the file name containing the proforma machine to use instead of the builtin proforma machine.
      -v, --version               version for genRandom

Use "genRandom [command] --help" for more information about a command.

## Commands

### binary

Generate a stream of (psudo)random ASCII '0' and '1' characters.

#### Usage

        genRandom binary [flags]

#### Flags

          --blocks=N string   Write N blocks. (default "1")
          --bs=BYTES string   Write up to BYTES bytes at a time. (default "512")
      -h, --help              help for binary

#### Global Flags

            --config string         config file (default is $HOME/.tnt2.yaml)
            --count=N string        starting block count.
        -o, --outputFile string     Name of the file containing the generated (psudo)random data. (default "-")
        -f, --proformafile string   the file name containing the proforma machine to use instead of the builtin proforma machine.

### completion

        Generate the autocompletion script for genRandom for the specified shell.
        See each sub-command's help for details on how to use the generated script.

#### Usage

        genRandom completion [command]

#### Available Commands

        bash        Generate the autocompletion script for bash
        fish        Generate the autocompletion script for fish
        powershell  Generate the autocompletion script for powershell
        zsh         Generate the autocompletion script for zsh

#### Flags

        -h, --help   help for completion

#### Global Flags

            --config string         config file (default is $HOME/.tnt2.yaml)
            --count=N string        starting block count.
        -o, --outputFile string     Name of the file containing the generated (psudo)random data. (default "-")
        -f, --proformafile string   the file name containing the proforma machine to use instead of the builtin proforma machine.

Use "genRandom completion [command] --help" for more information about a command.

### data

        Generate a stream of (psudo)random bytes.

#### Usage

        genRandom data [flags]

#### Flags

           --blocks string     Write N blocks. (default "1")
                          N and BYTES may be followed by the following multiplicative suffixes: c=1, 
                          w=2, b=512, kB=1000, K=1024, MB=1000*1000, M=1024*1024, GB=1000*1000*1000, 
                          G=1024*1024*1024, and so on for T, P, E, Z, Y.
            --bs=BYTES string   Write up to BYTES bytes at a time. (default "512")
        -h, --help              help for data

#### Global Flags

            --config string         config file (default is $HOME/.tnt2.yaml)
            --count=N string        starting block count.
        -o, --outputFile string     Name of the file containing the generated (psudo)random data. (default "-")
        -f, --proformafile string   the file name containing the proforma machine to use instead of the builtin proforma machine.

### Help

        Help provides help for any command in the application.
        Simply type genRandom help [path to command] for full details.

#### Usage

        genRandom help [command] [flags]

#### Flags

        -h, --help   help for help

#### Global Flags

            --config string         config file (default is $HOME/.tnt2.yaml)
            --count=N string        starting block count.
        -o, --outputFile string     Name of the file containing the generated (psudo)random data. (default "-")
        -f, --proformafile string   the file name containing the proforma machine to use instead of the builtin proforma machine.

### hex

        Generate a stream of (psudo)random hexadecimal encoded bytes.

#### Usage

        genRandom hex [flags]

#### Flags

            --blocks string     Write N blocks. (default "1")
                          N and BYTES may be followed by the following multiplicative suffixes: c=1, 
                          w=2, b=512, kB=1000, K=1024, MB=1000*1000, M=1024*1024, GB=1000*1000*1000, 
                          G=1024*1024*1024, and so on for T, P, E, Z, Y.
            --bs=BYTES string   Write up to BYTES bytes at a time. (default "512")
        -h, --help              help for hex

#### Global Flags

            --config string         config file (default is $HOME/.tnt2.yaml)
            --count=N string        starting block count.
        -o, --outputFile string     Name of the file containing the generated (psudo)random data. (default "-")
        -f, --proformafile string   the file name containing the proforma machine to use instead of the builtin proforma machine.

### points

        Generate a series of (psudo)random X,Y coordinates in a rectanglee of a given size

#### Usage

        genRandom points [flags]

#### Flags

        -g, --geometry string   The geometry of the rectangle in which to generate points expressed as WxH (eg. 640x480).
                                If geometry is present, it overide the height and width options.
        -y, --height int        Height of the rectangle in which to generate points. (default 480)
        -h, --help              help for points
        -n, --points=N string   Count of points to generate.
        -x, --width int         Width of the rectangle in which to generate points. (default 640)

#### Global Flags

            --config string         config file (default is $HOME/.tnt2.yaml)
            --count=N string        starting block count.
        -o, --outputFile string     Name of the file containing the generated (psudo)random data. (default "-")
        -f, --proformafile string   the file name containing the proforma machine to use instead of the builtin proforma machine.

### help

        Display version and detailed build information for tnt2.

#### Usage

        genRandom version [flags]

#### Flags

        -h, --help   help for version

#### Global Flags

            --config string         config file (default is $HOME/.tnt2.yaml)
            --count=N string        starting block count.
        -o, --outputFile string     Name of the file containing the generated (psudo)random data. (default "-")
        -f, --proformafile string   the file name containing the proforma machine to use instead of the builtin proforma machine.

---
## Notes

> N and BYTES may be followed by the following multiplicative suffixes: c=1, w=2, b=512, kB=1000, K=1024, MB=1000\*1000, M=1024\*1024, GB=1000\*1000\*1000, G=1024\*1024\*1024, and so on for T, P, E, Z, Y.

> The count can be a number (possibly followed by a multiplicative suffixes) or a fraction such as "1/2", "2/3", or "3/4".  If it is a fraction, then the starting block count is calculated by multiplying the maximal blocks generated by the tntEngine by the fraction.  Supplying a count will overide the stored count in the .genRand file, allowing for a repeatable stream of psuedo random data by giving the same secret key and starting block number.
