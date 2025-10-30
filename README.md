image-sampler
=

<pre>
iso_sampler is a command-line tool that allows users to sample files from ISO images based on specified criteria.

Usage:
  iso_sampler [command]

Available Commands:
  completion        Generate the autocompletion script for the specified shell
  help              Help about any command
  process-directory Process a directory
  sample-image      Sample files from an ISO image

Flags:
  -h, --help   help for iso_sampler
</pre>

process-directory
-
<pre>
Process a directory of images by sampling files based on specified criteria.

Usage:
  iso_sampler process-directory [flags]

Flags:
  -d, --directory string   Path to the directory containing ISO images
  -h, --help               help for process-directory
  -l, --limit int          Maximum number of directories to sample from per image (default 10)
  -o, --out string         Location to export sampled files (default "exports")
  </pre>

  sample-image
  -
  <pre>
  Sample files from an ISO image based on specified criteria such as directory limit and export location.

Usage:
  iso_sampler sample-image [flags]

Flags:
  -h, --help           help for sample-image
  -i, --image string   Path to the ISO image
  -l, --limit int      Maximum number of directories to sample from (default 10)
  -o, --out string     Location to export sampled files (default "exports")
  </pre>
