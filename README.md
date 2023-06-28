# Central Tools Standalone (CTstandalone)

This is a collection of standalone tools for performing various configuration and validation tasks via the Central API.  It consists of the following tools currently:

-------------------------------------
CTadd_var Version: 1.5
Author: Michael Gresham (michael.gresham@hpe.com)
-------------------------------------
usage: CTadd_var [-h|--help] [--token "<value>"] [--url "<value>"] [--initDB]
                 [--infile "<value>"] [--variable "<value>"] [--value
                 "<value>"] [-t|--test] [-c|--config "<value>"]

                 CTadd_var: Add/Update a variable in Central for a list of
                 device serial numbers.

Arguments:

  -h  --help      Print help information
      --token     Central API token if not using encrypted storage.
      --url       Central API URL if not using encrypted storage.
      --initDB    Initialize secure storage
      --infile    Input file consisting of a single device serial on each line
      --variable  Variable to create/update
      --value     Value to assign to the variable
  -t  --test      Enable test mode. No variables will be changed
  -c  --config    Config file location
  
-------------------------------------
CTarchive Version: 1.0
Author: Michael Gresham (michael.gresham@hpe.com)
-------------------------------------
usage: CTarchive [-h|--help] [--token "<value>"] [--url "<value>"] [--initDB]
                 [-t|--test] [--infile "<value>"] [-c|--config "<value>"]

                 CTarchive: Move serials to archive in Central.

Arguments:

  -h  --help    Print help information
      --token   Central API token if not using encrypted storage.
      --url     Central API URL if not using encrypted storage.
      --initDB  Initialize secure storage
  -t  --test    Enable test mode. No variables will be changed
      --infile  Input file consisting of a single device serial on each line
  -c  --config  Config file location

 
-------------------------------------
CTauto_commit Version: 1.0
Author: Michael Gresham (michael.gresham@hpe.com)
-------------------------------------
usage: CTauto_commit [-h|--help] [--token "<value>"] [--url "<value>"]
                     [--initDB] [--infile "<value>"] [--state "<value>"]
                     [-t|--test] [-c|--config "<value>"]

                     CTauto_commit: Enable/Disable autocommit for a list of
                     device serial numbers.

Arguments:

  -h  --help    Print help information
      --token   Central API token if not using encrypted storage.
      --url     Central API URL if not using encrypted storage.
      --initDB  Initialize secure storage
      --infile  Input file consisting of a single device serial on each line
      --state   Autocommit state: enable or disable
  -t  --test    Enable test mode. No variables will be changed
  -c  --config  Config file location

-------------------------------------
CTcentral_check Version: 1.0
Author: Michael Gresham (michael.gresham@ghpe.com)
-------------------------------------
usage: CTcentral_check [-h|--help] [--initDB] [-c|--config "<value>"]

                       CTcentral_check: Example program to access Central using
                       the API.

Arguments:

  -h  --help    Print help information
      --initDB  Initialize secure storage
  -c  --config  Config file location

-------------------------------------
CTdel_var Version: 1.5
Author: Michael Gresham (michael.gresham@hpe.com)
-------------------------------------
usage: CTdel_var [-h|--help] [--token "<value>"] [--url "<value>"] [--initDB]
                 [--infile "<value>"] [--variable "<value>"] [-t|--test]
                 [-c|--config "<value>"]

                 CTdel_var: Delete a variable in Central for a list of device
                 serial numbers.

Arguments:

  -h  --help      Print help information
      --token     Central API token if not using encrypted storage.
      --url       Central API URL if not using encrypted storage.
      --initDB    Initialize secure storage
      --infile    Input file consisting of a single device serial on each line
      --variable  Variable to delete
  -t  --test      Enable test mode. No variables will be changed
  -c  --config    Config file location


-------------------------------------
CTsite_serial Version: 1.5
Author: Michael Gresham (michael.gresham@hpe.com)
-------------------------------------
usage: CTsite_serials [-h|--help] [--token "<value>"] [--url "<value>"]
                      [--initDB] [--site "<value>"] [--dev_type "<value>"]
                      [-c|--config "<value>"]

                      CTsite_serial: Retrieve a list of device serial numbers
                      for a given site.

Arguments:

  -h  --help      Print help information
      --token     Central API token if not using encrypted storage.
      --url       Central API URL if not using encrypted storage.
      --initDB    Initialize secure storage
      --site      Target site name
      --dev_type  Device type requested: ALL|AP|SW|GW|MC
  -c  --config    Config file location


-------------------------------------
CTtemplate_add Version: 1.1
Author: Michael Gresham (michael.gresham@hpe.com)
-------------------------------------
usage: CTtemplate_add [-h|--help] [--token "<value>"] [--url "<value>"]
                      [--initDB] [-c|--config "<value>"]

                      CTtemplate_add: Add one or more template groups in
                      Central.

Arguments:

  -h  --help    Print help information
      --token   Central API token if not using encrypted storage.
      --url     Central API URL if not using encrypted storage.
      --initDB  Initialize secure storage
  -c  --config  Config file location

-------------------------------------
CTunarchive Version: 1.0
Author: Michael Gresham (michael.gresham@hpe.com)
-------------------------------------
usage: CTunarchive [-h|--help] [--token "<value>"] [--url "<value>"] [--initDB]
                   [-t|--test] [--infile "<value>"] [-c|--config "<value>"]

                   CTunarchive: Move serials from the archive to the live
                   inventory in Central.

Arguments:

  -h  --help    Print help information
      --token   Central API token if not using encrypted storage.
      --url     Central API URL if not using encrypted storage.
      --initDB  Initialize secure storage
  -t  --test    Enable test mode. No variables will be changed
      --infile  Input file consisting of a single device serial on each line
  -c  --config  Config file location


## All tools can use a CLI provided token OR a secure storage file for access to Central.  Using the secure storage file requires a one time call to any of the tools to initialize the secure storage.

     CTcentral_check-64.exe --initDB --config <secure storage file location>

## Compiled versions of the tools are provided for Windows 32/64, Linux and macOS 
