package main

import (
	"bufio" /* buff */
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
)

const INT_MAX = int(^uint(0) >> 1)

type selpg_args struct {
	start_page  int
	end_page    int
	in_filename string
	print_dest  string
	page_len    int  /* default 72, can be overriden by "-l number" on command line */
	page_type   bool /* default is 'l'; 'l' for lines-delimited, 'f' for form-feed-delimited */
}

type sp_args selpg_args

var progname string /* program name, for error messages */

func main() {
	write_file()
	sa := sp_args{}
	progname = os.Args[0] /* arg at index 0 is the command name itself (selpg) */
	get_args(&sa)
	check_args(&sa)
	process_input(&sa)
}

/* write file */
func write_file() {
	fout, err := os.Create("input")
	defer fout.Close()
	if err != nil {
		log.Fatalln(err)
	}
	for i := 1; i <= 100; i++ {
		if i%50 == 0 {
			_, err = fmt.Fprintf(fout, "%v\f\n", i)
		} else {
			_, err = fmt.Fprintln(fout, i)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

/*Pflag processes the args*/
func get_args(sa *sp_args) {
	flag.IntVarP(&sa.start_page, "start_page", "s", -1, "start page")
	flag.IntVarP(&sa.end_page, "end_page", "e", -1, "end page")
	flag.IntVarP(&sa.page_len, "page_length", "l", 72, "page len")
	flag.StringVarP(&sa.print_dest, "print_dest", "d", "", "print destination")
	flag.BoolVarP(&sa.page_type, "page_type", "f", false, "page type")
	flag.Parse()
}

func check_args(sa *sp_args) {
	/* Not enough args, minimum command is "selpg -sstartpage -eend_page"  */
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "\n%s: not enough arguments\n", progname)
		flag.Usage()
		os.Exit(1)
	}
	/* handle mandatory args first */
	/* handle arg - start_page -s*/
	if (sa.start_page <= 0) || (sa.start_page > INT_MAX) {
		fmt.Fprintf(os.Stderr, "\n%s: invalid start page %d\n", progname, sa.start_page)
		flag.Usage()
		os.Exit(2)
	}
	/* handle arg - end page -e*/
	if (sa.end_page <= 0) || (sa.end_page > INT_MAX) || (sa.end_page < sa.start_page) {
		fmt.Fprintf(os.Stderr, "\n%s: invalid end page %d\n", progname, sa.end_page)
		flag.Usage()
		os.Exit(3)
	}

	/* now handle optional args */
	/* handle arg - page_len -l*/
	if (sa.page_len < 1) || (sa.page_len > INT_MAX) {
		fmt.Fprintf(os.Stderr, "\n%s: invalid page length %d\n", progname, sa.page_len)
		flag.Usage()
		os.Exit(4)
	}
	/* handle arg - page_type -f */
	if (true == sa.page_type) && (sa.page_len != 72) {
		fmt.Fprintf(os.Stderr, "\n%s: command -l and -f are exclusive\n", progname)
		flag.Usage()
		os.Exit(5)
	}
	/* arg - print_dest -d is elided */

	/* there is one more arg - in_filename*/
	if len(flag.Args()) > 0 {
		_, err := os.Stat(flag.Args()[0])
		if err != nil && os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "\n%s: input file %s doesn't exist\n", progname, sa.in_filename)
			os.Exit(6)
		}
		sa.in_filename = flag.Args()[0]
	} else {
		sa.in_filename = ""
	}

	page_type := "page length"
	if sa.page_type == true {
		page_type = "end sign \f"
	}
	fmt.Fprintf(os.Stdout, "\nstart_page: %v", sa.start_page)
	fmt.Fprintf(os.Stdout, "\nend_page: %v", sa.end_page)
	fmt.Fprintf(os.Stdout, "\ninput_file: %s", sa.in_filename)
	fmt.Fprintf(os.Stdout, "\npage_length: %d", sa.page_len)
	fmt.Fprintf(os.Stdout, "\npage_type: %s", page_type)
	fmt.Fprintf(os.Stdout, "\nprint_destination: %s\n", sa.print_dest)
}

func process_input(sa *sp_args) {
	/* set the input source */
	var fin *os.File
	if sa.in_filename == "" {
		fin = os.Stdin
	} else {
		var err error
		fin, err = os.Open(sa.in_filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n%s: could not open input file \"%s\"\n", progname, sa.in_filename)
			os.Exit(7)
		}
		defer fin.Close() /* delay this close */
	}

	/* set the output destination */
	var fout io.WriteCloser
	if len(sa.print_dest) == 0 {
		fout = os.Stdout
		process_print(fout, fin, *sa) /*call process_print to print pages*/
	} else {
		fout = create_pipe(sa.print_dest)
		process_print(fout, fin, *sa)
		defer fout.Close() /* delay this close */
	}
}

func create_pipe(dest string) io.WriteCloser {
	var cmd *exec.Cmd

	cmd = exec.Command("lp", "-d"+dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fout, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pipe")
		log.Fatal(err)
	}
	return fout
}

func process_print(fout interface{}, fin *os.File, sa sp_args) {
	line_ctr, page_ctr, next_page := 0, 1, 1
	buf := bufio.NewReader(fin)
	var err error
	for true {
		var line string
		if sa.page_type { /*page type f*/
			line, err = buf.ReadString('\f')
			next_page = page_ctr + 1
		} else {
			line, err = buf.ReadString('\n')
			line_ctr++
			if line_ctr > sa.page_len {
				page_ctr++
				next_page = page_ctr
				line_ctr = 1
			}
		}

		if err == io.EOF {
			//fmt.Fprintln(os.Stderr, "[EOF]")
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if page_ctr >= sa.start_page && page_ctr <= sa.end_page {
			if std_out, ok := fout.(*os.File); ok {
				_, err = fmt.Fprintf(std_out, "%s", line)
			} else if std_pipe, ok := fout.(io.WriteCloser); ok {
				_, err = std_pipe.Write([]byte(line))
			} else {
				fmt.Fprintf(os.Stderr, "\n[Error]: fout type error\n")
				os.Exit(9)
			}
			if err != nil {
				log.Fatal(err)
			}
		}
		page_ctr = next_page
	}

	if page_ctr < sa.start_page {
		fmt.Fprintf(os.Stderr, "\n%s: start_page(%d) greater than total pages (%d), no output written\n", progname, sa.start_page, page_ctr)
		os.Exit(10)
	} else if page_ctr < sa.end_page {
		fmt.Fprintf(os.Stderr, "\n%s: end_page (%d) greater than total pages (%d), less output than expected\n", progname, sa.end_page, page_ctr)
		os.Exit(10)
	}
}
