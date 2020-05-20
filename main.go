// я потратил на это дерьмо два часа своей жизни
// напомните мне в следующий раз использовать awk как советовала интуиция, оки?
// upd: или хотя бы nex

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	re "regexp"
	"strconv"
	"time"
)

const latexpreamble = `% Эта таблица была сгенерирована автоматически
% ХРУПКОЕ, РУКАМИ НЕ ТРОГАТЬ!
% Вот это должно быть в преамбуле:

%\usepackage{longtable, tabu, caption, setspace}
%\captionsetup[table]{singlelinecheck=false,justification=raggedright,position=top,%
%        format=plain,font={large,rm,onehalfspacing},labelsep=endash,indention=0pt,%
%        skip=-4pt}
%\setlength{\LTcapwidth}{0pt}
%\def\GTBcch{\multicolumn{1}{|c|}}
%\def\GTBcct{\multicolumn{1}{c|}}
%\def\GTBnr{\\\hline}

`

var filename = flag.String("f", "", "filename")
var descsw = flag.String("d", "var", "print table of [func]tions or [var]iables")

//https://regex101.com/r/om4MWI/1
var fnc = re.MustCompile(`//(?:/|\s)*(.*)\n(` +
	`[A-Za-z_]+[A-Za-z:<>0-9_ ]*\**\s+([A-Za-z_]+[A-Za-z0-9_]*)\(.*\))`)

//https://regex101.com/r/4f3xgT/5
var vari = re.MustCompile(`//(?:/|\s)*([i-]?[o-]?[m-]?)\s*(.*)\n\s*([A-Za-z_]` +
	`+[A-Za-z:<>0-9_ ]*)\s+(\**\s*[A-Za-z_]+[A-Za-z0-9_]*(?:\[.*\])*)(?: +=.*)?;`)

func main() {
	flag.Parse()

	file, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Print("Usage:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if descsw == nil {
		fmt.Print("Usage:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	rand.Seed(time.Now().Unix())
	tblref := strconv.FormatInt(rand.Int63n(310), 10)

	if *descsw == "var" {
		varmatch := vari.FindAllSubmatch(file, -1)

		{
			fmt.Println(latexpreamble + `\noindent\begin{longtabu}{|X[1,l]|X[1,l]|X[3,l]|X[2,l]|}
\caption{Описание используемых переменных}\\\hline
\GTBcch{\bf Имя} & \GTBcct{\bf Тип} & \GTBcct{\bf Описание} & \GTBcct{\bf Направление}
\GTBnr\endfirsthead
\caption*{Продолжение таблицы \ref{ftbl` + tblref /* блять */ + `}}\\\hline%
\GTBcch{\bf Имя} & \GTBcct{\bf Тип} & \GTBcct{\bf Описание} & \cct{\bf Направление} 
\GTBnr\endhead`)
		}

		for _, match := range varmatch {
			fmt.Print(string(match[4]))
			fmt.Print(` & `)
			fmt.Print(string(match[3]))
			fmt.Print(` & `)
			fmt.Print(string(match[2]))
			fmt.Print(` & `)
			if bytes.ContainsRune(match[1], 'i') {
				fmt.Print(`входная, `)
			}
			if bytes.ContainsRune(match[1], 'o') {
				fmt.Print(`выходная, `)
			}
			if bytes.ContainsRune(match[1], 'm') {
				fmt.Print(`промежуточная`)
			}
			fmt.Println(`\GTBnr`)
		}
		fmt.Println(`\label{ftbl` + tblref + `}\GTBnr\end{longtabu}`)

	} else if *descsw == "func" {
		fncmatch := fnc.FindAllSubmatch(file, -1)

		{
			fmt.Println(latexpreamble + `\noindent\begin{longtabu}{|X[3,l]|X[4,l]|}
\caption{Функции, обеспечивающие работу программы}\\\hline
\GTBcch{\bf Имя} & \GTBcct{\bf Описание}\GTBnr\endfirsthead
\caption*{Продолжение таблицы \ref{dtbl` + tblref /* ты понял */ + `}}\\\hline
\GTBcch{\bf Имя} & \GTBcct{\bf Описание}\GTBnr\endhead`)
		}

		for _, match := range fncmatch {
			fmt.Print(`\tt `)
			fmt.Print(string(match[2]))
			fmt.Print(` & `)
			fmt.Print(string(match[1]))
			fmt.Println(`\GTBnr`)
		}
		fmt.Println(`\label{dtbl` + tblref + `}\GTBnr\end{longtabu}`)
	}
}
