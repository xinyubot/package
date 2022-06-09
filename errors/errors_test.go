package errors

import (
	"errors"
	"reflect"
	"testing"
)

func TestError_GetCode(t *testing.T) {
	type fields struct {
		code    int
		reason  string
		message string
		data    any
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{name: "case1", fields: fields{code: 500, reason: "", message: "", data: nil}, want: 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Error{
				code:    tt.fields.code,
				reason:  tt.fields.reason,
				message: tt.fields.message,
				data:    tt.fields.data,
			}
			if got := x.GetCode(); got != tt.want {
				t.Errorf("Error.GetCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_GetReason(t *testing.T) {
	type fields struct {
		code    int
		reason  string
		message string
		data    any
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "case1", fields: fields{code: 500, reason: "ASD", message: "", data: nil}, want: "ASD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Error{
				code:    tt.fields.code,
				reason:  tt.fields.reason,
				message: tt.fields.message,
				data:    tt.fields.data,
			}
			if got := x.GetReason(); got != tt.want {
				t.Errorf("Error.GetReason() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_GetMessage(t *testing.T) {
	type fields struct {
		code    int
		reason  string
		message string
		data    any
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "case1", fields: fields{code: 500, reason: "", message: "ASD", data: nil}, want: "ASD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Error{
				code:    tt.fields.code,
				reason:  tt.fields.reason,
				message: tt.fields.message,
				data:    tt.fields.data,
			}
			if got := x.GetMessage(); got != tt.want {
				t.Errorf("Error.GetMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_GetData(t *testing.T) {
	type fields struct {
		code    int
		reason  string
		message string
		data    any
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{name: "case1", fields: fields{code: 500, reason: "", message: "", data: map[string]int{"1": 1, "2": 2}}, want: map[string]int{"1": 1, "2": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Error{
				code:    tt.fields.code,
				reason:  tt.fields.reason,
				message: tt.fields.message,
				data:    tt.fields.data,
			}
			if got := x.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		code    int
		reason  string
		message string
		data    any
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "case1", fields: fields{code: 500, reason: "ASD", message: "", data: nil}, want: `code: 500, reason: ASD, message: , stack: , file: , data: <nil>, extraDataMap: map[string]interface {}(nil)`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				code:    tt.fields.code,
				reason:  tt.fields.reason,
				message: tt.fields.message,
				data:    tt.fields.data,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Is(t *testing.T) {
	type fields struct {
		code    int
		reason  string
		message string
		data    any
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "case1", fields: fields{code: 500, reason: "ASD", message: "", data: nil}, args: args{err: &Error{500, "ASD", "", "", "", nil, nil}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				code:    tt.fields.code,
				reason:  tt.fields.reason,
				message: tt.fields.message,
				data:    tt.fields.data,
			}
			if got := e.Is(tt.args.err); got != tt.want {
				t.Errorf("Error.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Equal(t *testing.T) {
	type fields struct {
		code    int
		reason  string
		message string
		data    any
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "case1", fields: fields{code: 500, reason: "ASD", message: "zxc", data: nil}, args: args{err: &Error{500, "ASD", "zxc", "", "", nil, nil}}, want: true},
		{name: "case2", fields: fields{code: 500, reason: "", message: "zxc", data: nil}, args: args{err: &Error{500, "ASD", "zxc", "", "", nil, nil}}, want: true},
		{name: "case3", fields: fields{code: 500, reason: "", message: "", data: nil}, args: args{err: &Error{500, "ASD", "zxc", "", "", nil, nil}}, want: false},
		{name: "case4", fields: fields{code: 501, reason: "ASD", message: "zxc", data: nil}, args: args{err: &Error{500, "ASD", "zxc", "", "", nil, nil}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				code:    tt.fields.code,
				reason:  tt.fields.reason,
				message: tt.fields.message,
				data:    tt.fields.data,
			}
			if got := e.Equal(tt.args.err); got != tt.want {
				t.Errorf("Error.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_DeepEqual(t *testing.T) {
	type fields struct {
		code    int
		reason  string
		message string
		data    any
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "case1", fields: fields{code: 500, reason: "ASD", message: "zxc", data: nil}, args: args{err: &Error{500, "ASD", "zxc", "", "", nil, nil}}, want: true},
		{name: "case2", fields: fields{code: 501, reason: "ASD", message: "zxc", data: nil}, args: args{err: &Error{500, "ASD", "zxc", "", "", nil, nil}}, want: false},
		{name: "case3", fields: fields{code: 500, reason: "", message: "zxc", data: nil}, args: args{err: &Error{500, "ASD", "zxc", "", "", nil, nil}}, want: false},
		{name: "case4", fields: fields{code: 500, reason: "ASD", message: "zxc", data: 1234}, args: args{err: &Error{500, "ASD", "zxc", "", "", nil, nil}}, want: false},
		{name: "case5", fields: fields{code: 500, reason: "ASD", message: "zxc", data: 12345}, args: args{err: &Error{500, "ASD", "zxc", "", "", 12345, nil}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				code:    tt.fields.code,
				reason:  tt.fields.reason,
				message: tt.fields.message,
				data:    tt.fields.data,
			}
			if got := e.DeepEqual(tt.args.err); got != tt.want {
				t.Errorf("Error.DeepEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		code    int
		reason  string
		message string
		data    any
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{name: "case1", args: args{code: 500, reason: "ASD", message: "", data: 123}, want: &Error{500, "ASD", "", "", "", 123, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.code, tt.args.reason, tt.args.message, tt.args.data); !got.Equal(tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{name: "case1", args: args{err: &Error{code: 500, reason: "ASD", message: "", data: 123}}, want: &Error{500, "ASD", "", "", "", 123, nil}},
		{name: "case2", args: args{err: errors.New("asd")}, want: &Error{500, "", "asd", "", "", []struct{}{}, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromError(tt.args.err); !got.Equal(tt.want) {
				t.Errorf("FromError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromErrorPro(t *testing.T) {
	type args struct {
		err    error
		code   int
		reason string
		data   any
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{name: "case1", args: args{err: &Error{code: 500, reason: "ASD", message: "", data: 123}, code: 500, reason: "ASD", data: 123}, want: &Error{500, "ASD", "", "", "", 123, nil}},
		{name: "case2", args: args{err: errors.New("asd"), code: 500, reason: "ASD", data: nil}, want: &Error{500, "ASD", "asd", "", "", []struct{}{}, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromErrorPro(tt.args.err, tt.args.code, tt.args.reason, tt.args.data); !got.Equal(tt.want) {
				t.Errorf("FromError() = %v, want %v", got, tt.want)
			}
		})
	}
}
