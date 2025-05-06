package server

import "testing"

func TestExtractThinking(t *testing.T) {
	tests := []struct {
		in, wantContent, wantThink string
	}{
		{
			in:          "hello <think> internal </think> world",
			wantContent: "hello  world",
			wantThink:   " internal ",
		},
		{
			in:          "<think>a</think><think>b</think>c",
			wantContent: "c",
			wantThink:   "a\nb",
		},
		{
			in:          "no think",
			wantContent: "no think",
			wantThink:   "",
		},
	}
	for i, tt := range tests {
		gotContent, gotThink := extractThinking(tt.in)
		if gotContent != tt.wantContent || gotThink != tt.wantThink {
			t.Errorf("case %d: got (%q,%q), want (%q,%q)", i, gotContent, gotThink, tt.wantContent, tt.wantThink)
		}
	}
}

func TestThinkingStreaming(t *testing.T) {

	type step struct {
		input          string
		wantThinking   string
		wantContent    string
		wantStateAfter thinkingParseState
	}

	cases := []struct {
		desc  string
		skip  bool
		steps []step
	}{
		{
			desc: "content without a thinking tag",
			steps: []step{
				{
					input:          "  abc",
					wantThinking:   "",
					wantContent:    "  abc",
					wantStateAfter: thinkingParseState_ThinkingDone,
				},
			},
		},
		{
			desc: "content before a thinking tag nerfs the thinking tag",
			steps: []step{
				{
					input:          "  abc <think>def</think> ghi",
					wantThinking:   "",
					wantContent:    "  abc <think>def</think> ghi",
					wantStateAfter: thinkingParseState_ThinkingDone,
				},
			},
		},
		{
			desc: "building up a thinking tag partially",
			// skip: true,
			steps: []step{
				{
					input:          "  <th",
					wantThinking:   "",
					wantContent:    "",
					wantStateAfter: thinkingParseState_LookingForOpening,
				},
				{
					input:          "in",
					wantThinking:   "",
					wantContent:    "",
					wantStateAfter: thinkingParseState_LookingForOpening,
				},
				{
					input:          "k>a",
					wantThinking:   "a",
					wantContent:    "",
					wantStateAfter: thinkingParseState_Thinking,
				},
			},
		},
		{
			desc: "partial closing tag",
			steps: []step{
				{
					input:          "<think>abc</th",
					wantThinking:   "abc",
					wantContent:    "",
					wantStateAfter: thinkingParseState_Thinking,
				},
				{
					input:          "ink>def",
					wantThinking:   "",
					wantContent:    "def",
					wantStateAfter: thinkingParseState_ThinkingDone,
				},
			},
		},
		{
			desc: "partial closing tag fakeout",
			steps: []step{
				{
					input:          "<think>abc</th",
					wantThinking:   "abc",
					wantContent:    "",
					wantStateAfter: thinkingParseState_Thinking,
				},
				{
					input:          "ing>def",
					wantThinking:   "</thing>def",
					wantContent:    "",
					wantStateAfter: thinkingParseState_ThinkingDone,
				},
				{
					input:          "ghi</thi",
					wantThinking:   "ghi",
					wantContent:    "",
					wantStateAfter: thinkingParseState_ThinkingDone,
				},
				{
					input:          "nk>jkl",
					wantThinking:   "",
					wantContent:    "jkl",
					wantStateAfter: thinkingParseState_ThinkingDone,
				},
			},
		},
	}

	for _, c := range cases {
		parser := thinkingParser{
			openingTag: "<think>",
			closingTag: "</think>",
		}
		if c.skip {
			continue
		}
		for i, step := range c.steps {
			thinking, content := parser.addContent(step.input)
			if content != step.wantContent || thinking != step.wantThinking {
				t.Errorf("case %q (step %d): got (%q,%q), want (%q,%q)", c.desc, i, content, thinking, step.wantContent, step.wantThinking)
			}
		}
	}
}
