# DMARC report CLI (pretty print DMARC reports)

This is a command-line tool to parse DMARC aggregate reports and output the
results in an easy to read manner.

This tools does not download the DMARC reports from a mailbox, but consumes the
previously downloaded reports from a local directory. This tool does not
connect to the Internet.

Given a directory, it will parse all reports for a given domain (from `.xml`,
`.xml.gz`, `.zip` files) and output a simplified DMARC report summary.

For example:

    $ dmarc-report ~/Downloads example.com

    # DKIM report for domain: example.com

    ## Reporter: google.com

    Begin                   End                     Policy: DKIM    SPF     Auth:   DKIM    SPF
    2023-09-22 02:00:00     2023-09-24 01:59:59             pass    pass            pass    pass
    2023-10-06 02:00:00     2023-10-07 01:59:59             pass    pass            pass    pass
    2023-10-29 02:00:00     2023-10-30 00:59:59             pass    pass            pass    pass
    2023-11-13 01:00:00     2023-11-14 00:59:59             pass    pass            pass    pass
    2023-12-04 01:00:00     2023-12-05 00:59:59             FAIL    pass            FAIL    pass
    2023-12-07 01:00:00     2023-12-08 00:59:59             pass    pass            pass    pass

    ## Reporter: Outlook.com

    Begin                   End                     Policy: DKIM    SPF     Auth:   DKIM    SPF
    2023-09-28 02:00:00     2023-09-29 02:00:00             pass    pass            pass    pass
    2023-10-03 02:00:00     2023-10-04 02:00:00             pass    pass            pass    pass
    2023-10-06 02:00:00     2023-10-07 02:00:00             pass    pass            pass    pass
    2023-11-12 01:00:00     2023-11-14 01:00:00             pass    pass            pass    pass
    2023-12-07 01:00:00     2023-12-08 01:00:00             pass    pass            pass    pass

## Installation

    go install github.com/candidtim/dmarc-report@1.0.0

## Usage

Usage:

    dmarc-report DIRECTORY DOMAIN

Example:

    dmarc-report ~/Downloads example.com

Latter will look for all files matching the DMARC report file name format
(`receiver "!" policy-domain "!" begin-timestamp "!" end-timestamp "." extension`),
decompress `.gz` and `.zip` files if necessary, and output an aggregated
summary.

## Features and limitations

 - Parses DMARC aggregate reports only (what is normally received by email when
   DMARC is set up for a domain)
 - Ignores duplicate reports (if there are multiple copies of the same report
   in the same directory; typically the case if some reports were extracted
   manually after they have been downloaded)
 - Merges subsequent reports (where the start time of a subsequent report is
   the same as the end time of a previous one)
 - Tested on most popular DMARC reports (Google, Outlook). Not extensively
   tested in various complex scenarios.


Note that this tool is originally implemented for individuals who own their domains,
set up DMARC, and want to review the reports regularly. Some more complex
scenarios may not be supported. However, feel free to submit issues, propose
features, or better yet - pull requests.
