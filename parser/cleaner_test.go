package parser_test

import (
	"testing"

	"github.com/sasalatart/batcoms/parser"
)

func TestClean(t *testing.T) {
	cc := []struct {
		input    string
		expected string
	}{
		{
			input:    "23 August 1942 – 2 February 1943[Note 1](5 months, 1 week and 3 days)",
			expected: "23 August 1942 – 2 February 1943 (5 months, 1 week and 3 days)",
		},
		{
			input:    "Stalingrad, Russian SFSR, Soviet Union(now Volgograd, Russia)",
			expected: "Stalingrad, Russian SFSR, Soviet Union (now Volgograd, Russia)",
		},
		{
			input:    "Soviet victory:[1]\n\nDestruction of the German 6th Army",
			expected: "Soviet victory: Destruction of the German 6th Army",
		},
		{
			input:    "\nInitial:270,000 personnel3,000 artillery pieces500 tanks600 aircraft, 1,600 by mid-September (Luftflotte 4)[Note 5][2]\nAt the time of the Soviet counter-offensive:c. 1,040,000 men[3][4]400,000+ Germans220,000 Italians200,000 Hungarians143,296 Romanians40,000 Hiwi10,250 artillery pieces500 tanks (140 Romanian)732 (402 operational) aircraft[5]:225[6]:87",
			expected: "Initial: 270,000 personnel 3,000 artillery pieces 500 tanks 600 aircraft, 1,600 by mid-September (Luftflotte 4). At the time of the Soviet counter-offensive: c. 1,040,000 men 400,000+ Germans 220,000 Italians 200,000 Hungarians 143,296 Romanians 40,000 Hiwi 10,250 artillery pieces 500 tanks (140 Romanian) 732 (402 operational) aircraft",
		},
		{
			input:    "\n1,129,619478,741 killed or missing650,878 wounded or sick[16]\n2,769 aircraft4,341 tanks (~150 by Romanians) (25-30% were total write-offs.[17])15,728 guns\nSee casualties section.",
			expected: "1,129,619. 478,741 killed or missing 650,878 wounded or sick. 2,769 aircraft 4,341 tanks (~150 by Romanians) (25-30% were total write-offs.) 15,728 guns. See casualties section.",
		},
		{
			input:    "\n8,300 killed,[1]3,400 captured",
			expected: "8,300 killed, 3,400 captured",
		},
		{
			input:    " .",
			expected: ".",
		},
		{
			input:    "word-1   word-2",
			expected: "word-1 word-2",
		},
		{
			input:    "Irregular whitespace. Another irregular whitespace.",
			expected: "Irregular whitespace. Another irregular whitespace.",
		},
	}

	for _, c := range cc {
		got := parser.Clean(c.input)
		if got != c.expected {
			t.Errorf("Expected\n%q\nto be parsed as:\n%q\nbut instead got:\n%q", c.input, c.expected, got)
		}
	}
}
