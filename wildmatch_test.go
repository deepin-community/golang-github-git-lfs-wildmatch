package wildmatch

import (
	"runtime"
	"testing"
)

type Case struct {
	Pattern   string
	Subject   string
	Match     bool
	Opts      []opt
	MatchOpts MatchOpts
}

func (c *Case) Assert(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			if c.Match {
				t.Errorf("could not parse: %s (%s)", c.Pattern, err)
			}
		}
	}()

	p := NewWildmatch(c.Pattern, c.Opts...)
	if (c.MatchOpts != MatchOpts{} && p.MatchWithOpts(c.Subject, c.MatchOpts) != c.Match) ||
		(c.MatchOpts == MatchOpts{} && p.Match(c.Subject) != c.Match) {
		if c.Match {
			t.Errorf("expected match: %s, %s", c.Pattern, c.Subject)
		} else {
			t.Errorf("unexpected match: %s, %s", c.Pattern, c.Subject)
		}
	}
}

var Cases = []*Case{
	{
		Pattern: `foo`,
		Subject: `foo`,
		Match:   true,
	},
	{
		Pattern: `bar`,
		Subject: `foo`,
		Match:   false,
	},
	{
		Pattern: `???`,
		Subject: `foo`,
		Match:   true,
	},
	{
		Pattern: `??`,
		Subject: `foo`,
		Match:   false,
	},
	{
		Pattern: `*`,
		Subject: `foo`,
		Match:   true,
	},
	{
		Pattern: `f*`,
		Subject: `foo`,
		Match:   true,
	},
	{
		Pattern: `*f`,
		Subject: `foo`,
		Match:   false,
	},
	{
		Pattern: `*foo*`,
		Subject: `foo`,
		Match:   true,
	},
	{
		Pattern: `*ob*a*r*`,
		Subject: `foobar`,
		Match:   true,
	},
	{
		Pattern: `*ab`,
		Subject: `aaaaaaabababab`,
		Match:   true,
	},
	{
		Pattern: `foo\*`,
		Subject: `foo*`,
		Match:   true,
	},
	{
		Pattern: `foo\*bar`,
		Subject: `foobar`,
		Match:   false,
	},
	{
		Pattern: `f\\oo`,
		Subject: `f\oo`,
		Match:   true,
	},
	{
		Pattern: `*[al]?`,
		Subject: `ball`,
		Match:   true,
	},
	{
		Pattern: `[ten]`,
		Subject: `ten`,
		Match:   false,
	},
	{
		Pattern: `**[!te]`,
		Subject: `ten`,
		Match:   true,
	},
	{
		Pattern: `**[!ten]`,
		Subject: `ten`,
		Match:   false,
	},
	{
		Pattern: `t[a-g]n`,
		Subject: `ten`,
		Match:   true,
	},
	{
		Pattern: `t[!a-g]n`,
		Subject: `ten`,
		Match:   false,
	},
	{
		Pattern: `t[!a-g]n`,
		Subject: `ton`,
		Match:   true,
	},
	{
		Pattern: `t[^a-g]n`,
		Subject: `ton`,
		Match:   true,
	},
	{
		Pattern: `]`,
		Subject: `]`,
		Match:   true,
	},
	{
		Pattern: `foo*bar`,
		Subject: `foo/baz/bar`,
		Match:   false,
	},
	{
		Pattern: `foo?bar`,
		Subject: `foo/bar`,
		Match:   false,
	},
	{
		Pattern: `foo[/]bar`,
		Subject: `foo/bar`,
		Match:   false,
	},
	{
		Pattern: `f[^eiu][^eiu][^eiu][^eiu][^eiu]r`,
		Subject: `foo/bar`,
		Match:   false,
	},
	{
		Pattern: `f[^eiu][^eiu][^eiu][^eiu][^eiu]r`,
		Subject: `foo-bar`,
		Match:   true,
	},
	{
		Pattern: `**/foo`,
		Subject: `foo`,
		Match:   true,
	},
	{
		Pattern: `**/foo`,
		Subject: `/foo`,
		Match:   true,
	},
	{
		Pattern: `**/foo`,
		Subject: `bar/baz/foo`,
		Match:   true,
	},
	{
		Pattern: `*/foo`,
		Subject: `bar/baz/foo`,
		Match:   false,
	},
	{
		Pattern: `**/bar*`,
		Subject: `foo/bar/baz`,
		Match:   false,
	},
	{
		Pattern: `**/bar/*`,
		Subject: `deep/foo/bar/baz`,
		Match:   true,
	},
	{
		Pattern: `**/bar/*`,
		Subject: `deep/foo/bar/baz/`,
		Match:   true,
	},
	{
		Pattern: `**/bar/**`,
		Subject: `deep/foo/bar/baz/`,
		Match:   true,
	},
	{
		Pattern: `**/bar/*`,
		Subject: `deep/foo/bar`,
		Match:   false,
	},
	{
		Pattern: `**/bar/**`,
		Subject: `deep/foo/bar/`,
		Match:   true,
	},
	{
		Pattern: `**/bar/**`,
		Subject: `deep/foo/bar`,
		Match:   false,
	},
	{
		Pattern: `**/bar/**/*`,
		Subject: `deep/foo/bar/`,
		Match:   true,
	},
	{
		Pattern: `**/bar/**/*`,
		Subject: `deep/foo/bar`,
		Match:   false,
	},
	{
		Pattern: `**/bar/**/*`,
		Subject: `deep/bar/bar`,
		Match:   false,
	},
	{
		Pattern: `*/bar/**`,
		Subject: `foo/bar/baz/x`,
		Match:   true,
	},
	{
		Pattern: `*/bar/**`,
		Subject: `deep/foo/bar/baz/x`,
		Match:   false,
	},
	{
		Pattern: `**/bar/*/*`,
		Subject: `deep/foo/bar/baz/x`,
		Match:   true,
	},
	{
		Pattern: `*.txt`,
		Subject: `foo/bar/baz.txt`,
		Match:   false,
	},
	{
		Pattern: `*.txt`,
		Subject: `你好-世界.txt`,
		Match:   true,
	},
	{
		Pattern: `你好-世界.txt`,
		Subject: `你好-世界.txt`,
		Match:   true,
	},
	{
		Pattern: `foo*`,
		Subject: `foobar`,
		Match:   true,
	},
	{
		Pattern: `*foo*`,
		Subject: `somethingfoobar`,
		Match:   true,
	},
	{
		Pattern: `*foo`,
		Subject: `barfoo`,
		Match:   true,
	},
	{
		Pattern: `a[c-c]st`,
		Subject: `acrt`,
		Match:   false,
	},
	{
		Pattern: `a[c-c]rt`,
		Subject: `acrt`,
		Match:   true,
	},
	{
		Pattern: `\`,
		Subject: `''`,
		Match:   false,
	},
	{
		Pattern: `\`,
		Subject: `\`,
		Match:   false,
	},
	{
		Pattern: `*/\`,
		Subject: `/\`,
		Match:   false,
	},
	{
		Pattern: `foo`,
		Subject: `foo`,
		Match:   true,
	},
	{
		Pattern: `@foo`,
		Subject: `@foo`,
		Match:   true,
	},
	{
		Pattern: `@foo`,
		Subject: `foo`,
		Match:   false,
	},
	{
		Pattern: `\[ab]`,
		Subject: `[ab]`,
		Match:   true,
	},
	{
		Pattern: `[[]ab]`,
		Subject: `[ab]`,
		Match:   true,
	},
	{
		Pattern: `[[:]ab]`,
		Subject: `[ab]`,
		Match:   true,
	},
	{
		Pattern: `[[::]ab]`,
		Subject: `[ab]`,
		Match:   false,
	},
	{
		Pattern: `[[:digit]ab]`,
		Subject: `[ab]`,
		Match:   false,
	},
	{
		Pattern: `[\[:]ab]`,
		Subject: `[ab]`,
		Match:   true,
	},
	{
		Pattern: `\??\?b`,
		Subject: `?a?b`,
		Match:   true,
	},
	{
		Pattern: `''`,
		Subject: `foo`,
		Match:   false,
	},
	{
		Pattern: `**/t[o]`,
		Subject: `foo/bar/baz/to`,
		Match:   true,
	},
	{
		Pattern: `[[:alpha:]][[:digit:]][[:upper:]]`,
		Subject: `a1B`,
		Match:   true,
	},
	{
		Pattern: `[[:digit:][:upper:][:space:]]`,
		Subject: `a`,
		Match:   false,
	},
	{
		Pattern: `[[:digit:][:upper:][:space:]]`,
		Subject: `A`,
		Match:   true,
	},
	{
		Pattern: `[[:digit:][:upper:][:space:]]`,
		Subject: `1`,
		Match:   true,
	},
	{
		Pattern: `[[:digit:][:upper:][:spaci:]]`,
		Subject: `1`,
		Match:   false,
	},
	{
		Pattern: `'`,
		Subject: `'`,
		Match:   true,
	},
	{
		Pattern: `[[:digit:][:upper:][:space:]]`,
		Subject: `.`,
		Match:   false,
	},
	{
		Pattern: `[[:digit:][:punct:][:space:]]`,
		Subject: `.`,
		Match:   true,
	},
	{
		Pattern: `[[:xdigit:]]`,
		Subject: `5`,
		Match:   true,
	},
	{
		Pattern: `[[:xdigit:]]`,
		Subject: `f`,
		Match:   true,
	},
	{
		Pattern: `[[:xdigit:]]`,
		Subject: `D`,
		Match:   true,
	},
	{
		Pattern: `[[:alnum:][:alpha:][:blank:][:cntrl:][:digit:][:graph:][:lower:][:print:][:punct:][:space:][:upper:][:xdigit:]]`,
		Subject: `_`,
		Match:   true,
	},
	{
		Pattern: `[^[:alnum:][:alpha:][:blank:][:cntrl:][:digit:][:lower:][:space:][:upper:][:xdigit:]]`,
		Subject: `.`,
		Match:   true,
	},
	{
		Pattern: `[a-c[:digit:]x-z]`,
		Subject: `5`,
		Match:   true,
	},
	{
		Pattern: `[a-c[:digit:]x-z]`,
		Subject: `b`,
		Match:   true,
	},
	{
		Pattern: `[a-c[:digit:]x-z]`,
		Subject: `y`,
		Match:   true,
	},
	{
		Pattern: `[a-c[:digit:]x-z]`,
		Subject: `q`,
		Match:   false,
	},
	{
		Pattern: `[\\-^]`,
		Subject: `]`,
		Match:   true,
	},
	{
		Pattern: `[\\-^]`,
		Subject: `[`,
		Match:   false,
	},
	{
		Pattern: `a[]b`,
		Subject: `ab`,
		Match:   false,
	},
	{
		Pattern: `a[]b`,
		Subject: `a[]b`,
		Match:   false,
	},
	{
		Pattern: `[!`,
		Subject: `ab`,
		Match:   false,
	},
	{
		Pattern: `[-`,
		Subject: `ab`,
		Match:   false,
	},
	{
		Pattern: `[-]`,
		Subject: `-`,
		Match:   true,
	},
	{
		Pattern: `[a-`,
		Subject: `-`,
		Match:   false,
	},
	{
		Pattern: `[!a-`,
		Subject: `-`,
		Match:   false,
	},
	{
		Pattern: `'`,
		Subject: `'`,
		Match:   true,
	},
	{
		Pattern: `'[`,
		Subject: `0`,
		Match:   false,
	},
	{
		Pattern: `[---]`,
		Subject: `-`,
		Match:   true,
	},
	{
		Pattern: `[------]`,
		Subject: `-`,
		Match:   true,
	},
	{
		Pattern: `[!------]`,
		Subject: `a`,
		Match:   true,
	},
	{
		Pattern: `[a^bc]`,
		Subject: `^`,
		Match:   true,
	},
	{
		Pattern: `[\]`,
		Subject: `\`,
		Match:   false,
	},
	{
		Pattern: `[\\]`,
		Subject: `\`,
		Match:   true,
	},
	{
		Pattern: `[!\\]`,
		Subject: `\`,
		Match:   false,
	},
	{
		Pattern: `[A-\\]`,
		Subject: `G`,
		Match:   true,
	},
	{
		Pattern: `b*a`,
		Subject: `aaabbb`,
		Match:   false,
	},
	{
		Pattern: `*ba*`,
		Subject: `aabcaa`,
		Match:   false,
	},
	{
		Pattern: `[,]`,
		Subject: `,`,
		Match:   true,
	},
	{
		Pattern: `[\\,]`,
		Subject: `,`,
		Match:   true,
	},
	{
		Pattern: `[\\,]`,
		Subject: `\`,
		Match:   true,
	},
	{
		Pattern: `[,-.]`,
		Subject: `-`,
		Match:   true,
	},
	{
		Pattern: `[,-.]`,
		Subject: `+`,
		Match:   false,
	},
	{
		Pattern: `[,-.]`,
		Subject: `-.]`,
		Match:   false,
	},
	{
		Pattern: `-*-*-*-*-*-*-12-*-*-*-m-*-*-*`,
		Subject: `-adobe-courier-bold-o-normal--12-120-75-75-m-70-iso8859-1`,
		Match:   true,
	},
	{
		Pattern: `-*-*-*-*-*-*-12-*-*-*-m-*-*-*`,
		Subject: `-adobe-courier-bold-o-normal--12-120-75-75-X-70-iso8859-1`,
		Match:   false,
	},
	{
		Pattern: `-*-*-*-*-*-*-12-*-*-*-m-*-*-*`,
		Subject: `-adobe-courier-bold-o-normal--12-120-75-75-/-70-iso8859-1`,
		Match:   false,
	},
	{
		Pattern: `**/*a*b*g*n*t`,
		Subject: `abcd/abcdefg/abcdefghijk/abcdefghijklmnop.txt`,
		Match:   true,
	},
	{
		Pattern: `**/*a*b*g*n*t`,
		Subject: `abcd/abcdefg/abcdefghijk/abcdefghijklmnop.txtz`,
		Match:   false,
	},
	{
		Pattern: `file[[:space:]]with[[:space:]]spaces.\#`,
		Subject: `file with spaces.#`,
		Match:   true,
	},
	{
		Pattern: `foo`,
		Subject: `FOO`,
		Match:   false,
	},
	{
		Pattern: `foo`,
		Subject: `FOO`,
		Opts:    []opt{CaseFold},
		Match:   true,
	},
	{
		Pattern: `**/a*.txt`,
		Subject: `foo-a.txt`,
		Match:   false,
	},
	{
		Pattern: `*.txt`,
		Subject: `file.txt`,
		Opts:    []opt{Basename},
		Match:   true,
	},
	{
		Pattern: `file.txt`,
		Subject: `file.txt`,
		Opts:    []opt{Basename, Contents},
		Match:   true,
	},
	{
		Pattern: `*.txt`,
		Subject: `path/to/file.txt`,
		Opts:    []opt{Basename},
		Match:   true,
	},
	{
		Pattern: `path/to/*.txt`,
		Subject: `path/to/file.txt`,
		Opts:    []opt{Basename},
		Match:   true,
	},
	{
		Pattern: `path/to/*.txt`,
		Subject: `path/to/file.txt`,
		Match:   true,
	},
	{
		Pattern: `path/to/*.txt`,
		Subject: `outside/of/path/to/file.txt`,
		Opts:    []opt{Basename},
		Match:   false,
	},
	{
		Pattern: `path/to/*.txt`,
		Subject: `path/to/some/intermediaries/to/file.txt`,
		Opts:    []opt{Basename},
		Match:   false,
	},
	{
		Pattern: `path/`,
		Subject: `path/to/some/intermediaries/to/file.txt`,
		Match:   false,
	},
	{
		// GitAttribute-style matching directory.
		// false becalse gitattribute never matches directories.
		Pattern:   `anotherfile.txt/`,
		Subject:   `anotherfile.txt`,
		Opts:      []opt{GitAttributes},
		MatchOpts: MatchOpts{IsDirectory: true},
		Match:     false,
	},
	{
		// gitAttribute-style matching normal file.
		// false as gitattribute matches ending in '/' indicate
		// trying to match directory but gitattribute never matches directory
		Pattern: `anotherfile1.txt/`,
		Subject: `anotherfile1.txt`,
		Opts:    []opt{GitAttributes},
		Match:   false,
	},
	{
		// gitignore-style matching directory.
		Pattern:   `anotherfile2.txt/`,
		Subject:   `anotherfile2.txt`,
		MatchOpts: MatchOpts{IsDirectory: true},
		Match:     true,
	},
	{
		Pattern: `anotherfile3.txt/`,
		Subject: `anotherfile3.txt`,
		Match:   false,
	},
	{
		Pattern: `anotherfile4.txt`,
		Subject: `anotherfile4.txt/`,
		Opts:    []opt{GitAttributes},
		Match:   false,
	},
	{
		Pattern: `**/pdfkit.frameworks/pdfkit/**`,
		Subject: `MyFolder/libs/pdfkit.frameworks/pdfkit`,
		Match:   false,
	},
	{
		Pattern:   `foo/`,
		Subject:   `bar/baz/foo`,
		MatchOpts: MatchOpts{IsDirectory: true},
		Match:     true,
	},
	{
		Pattern:   `foo/`,
		Subject:   `foo`,
		MatchOpts: MatchOpts{IsDirectory: true},
		Match:     true,
	},
	{
		Pattern: `foo/`,
		Subject: `foo/`,
		Match:   true,
	},
	{
		Pattern: `/foo/`,
		Subject: `foo/`,
		Match:   true,
	},
	{
		Pattern: `big/b`,
		Subject: `big/b/b1`,
		Opts:    []opt{Contents},
		Match:   true,
	},
	{
		Pattern: `big`,
		Subject: `big/b/b1`,
		Opts:    []opt{Contents},
		Match:   true,
	},
	{
		Pattern: `b`,
		Subject: `big/b/b1`,
		Opts:    []opt{Contents},
		Match:   true,
	},
	{
		Pattern: `/foo/`,
		Subject: `foo/`,
		Opts:    []opt{Contents},
		Match:   true,
	},
	{
		Pattern: `/foo/`,
		Subject: `foo/`,
		Opts:    []opt{Basename, Contents},
		Match:   true,
	},
	{
		Pattern: `/foo`,
		Subject: `foo`,
		Match:   true,
	},
	{
		Pattern: `/foo/filename.txt`,
		Subject: `foo/filename.txt`,
		Match:   true,
	},
	{
		Pattern: `/foo/filename.txt`,
		Subject: `bar/foo/filename.txt`,
		Match:   false,
	},
	{
		Pattern: `/foo/*.txt`,
		Subject: `foo/filename.txt`,
		Match:   true,
	},
	{
		Pattern: `/*.txt`,
		Subject: `foo/filename.txt`,
		Match:   false,
	},
	{
		Pattern: `/foo/*.txt`,
		Subject: `bar/foo/filename.txt`,
		Match:   false,
	},
	{
		Pattern:   `/foo/`,
		Subject:   `foo`,
		MatchOpts: MatchOpts{IsDirectory: true},
		Match:     true,
	},
	{
		Pattern: `/foo/`,
		Subject: `foo/filename.txt`,
		Opts:    []opt{Contents},
		Match:   true,
	},
	{
		Pattern: `/foo/**`,
		Subject: `foo/filename.txt`,
		Match:   true,
	},
	{
		Pattern: `path/`,
		Subject: `path/to/some/intermediaries/to/file.txt`,
		Opts:    []opt{Contents},
		Match:   true,
	},
	{
		Pattern: `to/`,
		Subject: `path/to/some/intermediaries/to/file.txt`,
		Opts:    []opt{Contents},
		Match:   true,
	},
	{
		Pattern: `nonexistent/`,
		Subject: `path/to/some/intermediaries/to/file.txt`,
		Opts:    []opt{Contents},
		Match:   false,
	},
}

func TestWildmatch(t *testing.T) {
	for _, c := range Cases {
		c.Assert(t)
	}
}

type SlashCase struct {
	Given  string
	Expect string
}

func (c *SlashCase) Assert(t *testing.T) {
	got := slashEscape(c.Given)

	if c.Expect != got {
		t.Errorf("wildmatch: expected slashEscape(\"%s\") -> %s, got: %s",
			c.Given,
			c.Expect,
			got,
		)
	}
}

func TestSlashEscape(t *testing.T) {
	for _, c := range []*SlashCase{
		{Given: ``, Expect: ``},
		{Given: `foo/bar`, Expect: `foo/bar`},
		{Given: `foo\bar`, Expect: `foo/bar`},
		{Given: `foo\*bar`, Expect: `foo\*bar`},
		{Given: `foo\?bar`, Expect: `foo\?bar`},
		{Given: `foo\[bar`, Expect: `foo\[bar`},
		{Given: `foo\]bar`, Expect: `foo\]bar`},
		{Given: `foo\#bar`, Expect: `foo\#bar`},
	} {
		c.Assert(t)
	}
}

func TestCaseFold(t *testing.T) {
	m := NewWildmatch("*.bin", SystemCase)
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		if !m.Match("UPCASE.BIN") {
			t.Errorf("wildmatch: expected system case to be folding")
		}
	} else if m.Match("UPCASE.BIN") {
		t.Errorf("wildmatch: expected system case to be non-folding")
	}
}
