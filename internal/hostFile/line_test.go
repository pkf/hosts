package hostfile

import (
	"reflect"
	"testing"
	"time"
)

func TestTransLine(t *testing.T) {
	type args struct {
		raw   string
		hosts map[string]string
	}
	var hosts = map[string]string{
		"live.github.com":    "140.82.114.26",
		"central.github.com": "140.82.112.21",
		"alive.github.com":   "140.82.113.25",
		"demo1.github.com":   "140.82.113.29",
		"demo2.github.com":   "140.82.113.77",
		"demo3.github.com":   "140.82.113.72",
		"demo4.github.com":   "140.82.113.71",
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "换IP",
			args: args{
				raw:   `127.0.0.1 demo1.github.com`,
				hosts: hosts,
			},
			want: []string{
				`140.82.113.29 demo1.github.com #` + time.Now().Format(time.RFC3339),
			},
		},
		{
			name: "多个换IP",
			args: args{
				raw:   `127.0.0.1 demo2.github.com central.github.com`,
				hosts: hosts,
			},
			want: []string{
				`140.82.113.77 demo2.github.com #` + time.Now().Format(time.RFC3339),
				`140.82.112.21 central.github.com #` + time.Now().Format(time.RFC3339),
			},
		},
		{
			name: "多个只换一个IP",
			args: args{
				raw:   `127.0.0.1 demo3.github.com notfound.github.com localhost`,
				hosts: hosts,
			},
			want: []string{
				`140.82.113.72 demo3.github.com #` + time.Now().Format(time.RFC3339),
				`127.0.0.1 notfound.github.com localhost`,
			},
		},
		{
			name: "多个只换前后IP",
			args: args{
				raw:   `127.0.0.1 live.github.com notfound.github.com demo4.github.com`,
				hosts: hosts,
			},
			want: []string{
				`140.82.114.26 live.github.com #` + time.Now().Format(time.RFC3339),
				`140.82.113.71 demo4.github.com #` + time.Now().Format(time.RFC3339),
				`127.0.0.1 notfound.github.com`,
			},
		},
		{
			name: "前面有被用过",
			args: args{
				raw:   `127.0.0.1 live.github.com notfound.github.com alive.github.com`,
				hosts: hosts,
			},
			want: []string{
				`140.82.113.25 alive.github.com #` + time.Now().Format(time.RFC3339),
				`127.0.0.1 notfound.github.com`,
			},
		},
		{
			name: "只有注释",
			args: args{
				raw:   `# test......`,
				hosts: hosts,
			},
			want: []string{
				`# test......`,
			},
		},
		{
			name: "带注释重新格式化",
			args: args{
				raw:   `127.0.0.1 abc.com   # test......`,
				hosts: hosts,
			},
			want: []string{
				`127.0.0.1 abc.com # test......`,
			},
		},
		{
			name: "多域名",
			args: args{
				raw:   `127.0.0.1 abc.com def.io   # test......`,
				hosts: hosts,
			},
			want: []string{
				`127.0.0.1 abc.com def.io # test......`,
			},
		},
		{
			name: "错误IP",
			args: args{
				raw:   `127.0.1   # test......`,
				hosts: hosts,
			},
			want: []string{
				`127.0.1   # test......`,
			},
		},
		{
			name: "空白行",
			args: args{
				raw: ` 	 `,
				hosts: hosts,
			},
			want: []string{
				` 	 `,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransLine(tt.args.raw, tt.args.hosts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
