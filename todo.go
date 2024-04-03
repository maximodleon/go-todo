package todo

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)


type item struct {
  Task string
  Done bool
  CreatedAt time.Time
  CompletedAt time.Time
}

type List []item

type Stringer interface {
  String() string
}

func (l *List) Add (task string) {
  t := item{
    Task: task,
    Done: false,
    CreatedAt: time.Now(),
    CompletedAt: time.Time{},
  }

  *l = append(*l, t)
}

func (l *List) Complete(i int) error {
  ls := *l

  if i <= 0 || i > len(ls) {
    return fmt.Errorf("Item %d doest not exist", i)
  }

  // We need to adjust to be 0 index based
  ls[i-1].Done = true
  ls[i-1].CompletedAt = time.Now()

  return nil
}

func (l *List) Delete(i int) error {
  ls := *l

  if i <=0 || i > len(ls) {
    return fmt.Errorf("Item %d does not exist", i)
  }

  *l = append(ls[:i-1],ls[i:]...)
  return nil
}

func (l *List) Save(filename string) error {
  js, err := json.Marshal(l)

  if err != nil {
    return err
  }

  return os.WriteFile(filename, js, 0644)
}

func (l *List) Get(filename string) error {
  file, err := os.ReadFile(filename)

  if err != nil {
    if errors.Is(err, os.ErrNotExist) {
      return nil
    }

    return err
  }

  if len(file) == 0 {
    return nil
  }

  return json.Unmarshal(file, l)
}

func (l *List) String() string {
  verbose := false
  skipDone := false
  flag.Visit(func(f *flag.Flag){
    if f.Name == "v" {
      verbose = true
    } else if f.Name == "skipDone" {
      skipDone = true
    }
  })

  formatted := ""

  for k,t := range *l {
     prefix := ""
     suffix := ""
     if t.Done {
      if skipDone {
        continue
      }
      prefix = "X "
    }

    if verbose {
      suffix = "\n    Created at: " + t.CreatedAt.String()
    }

    formatted += fmt.Sprintf("%s%d: %s%s\n", prefix, k+1, t.Task, suffix)
  }

  return formatted
}
