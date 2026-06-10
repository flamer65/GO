# 🏋️ Slices — Exercises & Questions

---

## Section 1: Conceptual Questions (Slice Internals)

### 1a. Slice Header

A slice header has 3 components. Name them and explain what each one stores.

```
Your answer:
1. pointer
2. length
3. capacity
```

### 1b. `make` vs Literal

What is the difference between these two?

```go
a := make([]int, 5)
b := make([]int, 0, 5)
```

```
Length of a: 5
Capacity of a: 5
What a looks like: [0,0,0,0,0]

Length of b: 0
Capacity of b: 5
What b looks like: []
```

### 1c. Zero Value

What is the zero value of a slice? Can you `append` to a nil slice?

```go
var s []int
fmt.Println(s == nil)    // ? true
fmt.Println(len(s))      // ? 0
s = append(s, 1)         // does this panic? no result [1]
fmt.Println(s)           // ? [1]
```

```
Your answers: true, 0, no, [1]
```

---

## Section 2: Predict the Output

### 2a. Shared Backing Array

```go
original := []int{10, 20, 30, 40, 50}
sub := original[1:3]

sub[0] = 999

fmt.Println(original) [10, 999, 30, 40, 50]
fmt.Println(sub) [99, 30]
```

```
original output: [10, 999, 30, 40, 50]
sub output: [99, 30]
Why: the sub-slices shares same backing array
```

### 2b. Append and Capacity

```go
s := make([]int, 0, 3)
s = append(s, 1)
s = append(s, 2)
s = append(s, 3)
fmt.Printf("len=%d, cap=%d\n", len(s), cap(s))

s = append(s, 4)
fmt.Printf("len=%d, cap=%d\n", len(s), cap(s))
```

```
First Printf: 3
Second Printf: 3
What happened on the 4th append: the capacity increase to 6
```

### 2c. Sub-slice Capacity

```go
s := []int{1, 2, 3, 4, 5, 6, 7, 8}
a := s[2:5]

fmt.Println(a)
fmt.Println(len(a))
fmt.Println(cap(a))
```

```
a value: [3,4,5]
len: 3
cap: 6
Why is cap that value: because share same back tracking not provided max to control capacity
```

### 2d. The Tricky Append

```go
s := []int{1, 2, 3, 4, 5}
a := s[1:3]

fmt.Println(a)       // ? [2, 3]
fmt.Println(cap(a))  // ? 4

a = append(a, 99)

fmt.Println(a)       // ? [2,3,99]
fmt.Println(s)       // did s change?
```

```
a before append: [2, 3]
cap of a:  4
a after append: [2,3,99]
s after append: [1, 2, 3, 99, 5]
Why did s change (or not): a shares the backing array AND had spare capacity, so append wrote 99 into s[3], overwriting the 4!
```

### 2e. Copy Behavior

```go
src := []int{1, 2, 3, 4, 5}
dst := make([]int, 3)

n := copy(dst, src)

fmt.Println(dst)
fmt.Println(n)
```

```
dst: [1, 2, 3]
n: 3
Why only 3 elements: copy looks for the shorter dst is shorter so onlt first 3 elements. copy doesn't goes destination slice capacity
```

### 2f. Append Creates New Array?

```go
s := []int{1, 2, 3}
s2 := append(s, 4)

s[0] = 999

fmt.Println(s)
fmt.Println(s2)
```

```
s: [1,2,3]
s2: [1,2,3,4]
Does s2 see the 999? Why: No, beacuse the s2 is has new address due to append
```

---

## Section 3: Slice Deletion & Insertion

### 3a. Delete at Index

Given `s := []int{10, 20, 30, 40, 50}`, write a one-liner to delete the element at index 2 (value 30).

```go
s = append(s[:2],s[3:]...)
```

What is the result?

```
s: [10,20,40,50]
```

### 3b. Insert at Index

Given `s := []int{1, 2, 4, 5}`, write code to insert the value `3` at index 2.

```go
s = append(s[:3],s[2:]....) s[2] = 3
```

What is the result? [1,2,3,4,5]

```
s:
```

### 3c. Memory Safety

After this deletion, what happens to the element that "fell off"?

```go
s := []*SomeStruct{a, b, c, d, e}
s = append(s[:2], s[3:]...) // delete index 2 ("c")
```

```
Is "c" garbage collected? Why or why not: d overrrides the c, c is garbage collected and where as there is memory leak on the last 4 index so
```

---

## Section 4: Three-Index Slice

### 4a. What Does the Third Index Do?

```go
s := []int{1, 2, 3, 4, 5, 6, 7, 8}
a := s[2:5:6]

fmt.Println(a)
fmt.Println(len(a))
fmt.Println(cap(a))
```

```
a: [3,4,5]
len:3
cap: 4
Why is cap different from s[2:5]: because the third puts cap on the capacity of the slice, that is max upto 6 index
```

### 4b. When Would You Use It?

Explain in your own words when the three-index slice `[low:high:max]` is useful.

```
Your answer: to prevent use of the shared-array as we sahring with the referrence so can put cap to the capacity of the sub-slice
```

---

## Section 5: Variadic and `...` Operator

### 5a. What Does `...` Do?

Explain the two uses of `...` in Go:

```go
// Use 1:
func sum(nums ...int) int { }

// Use 2:
s := []int{1, 2, 3}
result := append(s, []int{4, 5}...)
```

```
Use 1 meaning: collect multiple values while calling the function into single slice
Use 2 meaning: spread all the values of slice as indviual argument
```

### 5b. Predict the Output

```go
func add(s []int, vals ...int) []int {
    return append(s, vals...)
}

a := []int{1, 2}
b := add(a, 3, 4, 5)
fmt.Println(b)
```

```
Output: [1,2,3,4,5]
```

---

## Section 6: Slice Growth

### 6a. Growth Pattern

Starting from an empty slice, what capacity does Go allocate after each append?

```go
var s []int
for i := 0; i < 10; i++ {
    s = append(s, i)
    fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
}
```

Fill in the table:

| After append | len | cap |
| ------------ | --- | --- |
| 0            | 1   | 1   |
| 1            | 2   | 2   |
| 2            | 3   | 4   |
| 3            | 4   | 4   |
| 4            | 5   | 8   |
| 5            | 6   | 8   |
| 6            | 7   | 8   |
| 7            | 8   | 8   |
| 8            | 9   | 16  |
| 9            | 10  | 16  |

```
What is the general growth strategy: 2 * last capacity
```

---

## Section 7: Coding Exercises

### 7a. Flatten

Write a function that takes a `[][]int` and returns a single `[]int` with all elements.

```
Example: [[1,2],[3,4],[5]] → [1,2,3,4,5]
```

```go
func flatten(s [][]int) []int {
    arr := []int{}
  for _,num := range s{
    for _,j := range num{
        arr = append(arr,j)
    }
  }
  return arr
}
```

### 7b. Zip

Write a function that takes two slices and returns a slice of pairs (as [2]int arrays).

```
Example: zip([1,2,3], [4,5,6]) → [[1,4], [2,5], [3,6]]
If lengths differ, stop at the shorter one.
```

```go
func zip(a, b []int) [][2]int {
 n := len(a)
 if len(b) < n {
    n = len(b)
 }
 arr := [][2]int
 for i := 0; i < n ; i++{
    pair := [2]int{a[i], b[i]}
    arr = append(arr, pair)
 }
 return pair
}
```

### 7c. Sum & Average

Write two functions:

```go
func sum(s []int) int {
  s := 0;
  for _,num := range s{
    sum += num
  }
  return sum
}

func average(s []int) float64 {
   sum := sum(s)
   n := len(s)
  return float64(sum/n)
}
```

### 7d. Max & Min

Write a function that returns both the max and min of a slice.

```go
func minMax(s []int) (int, int) {

}
```

### 7e. Map (transform)

Write a function that applies a transform function to every element.

```
Example: mapSlice([1,2,3], func(n int) int { return n*n }) → [1,4,9]
```

```go
func mapSlice(s []int, fn func(int) int) []int {
  arr := []int{}
  for _,num := range s{
    n := fn(num)
    arr = append(n)
  }
  return arr
}
```

### 7f. Reduce

Write a function that reduces a slice to a single value using an accumulator.

```
Example: reduce([1,2,3,4], 0, func(acc, n int) int { return acc + n }) → 10
```

```go
func reduce(s []int, initial int, fn func(int, int) int) int {

   for _, num := range s{
      initial = fn(initial,num)
   }
   return initial
}
```

### 7g. Intersection

Write a function that returns elements present in BOTH slices.

```
Example: intersection([1,2,3,4], [3,4,5,6]) → [3,4]
```

```go
func intersection(a, b []int) []int {
   hash := make(map[int]bool)
   for _,num := range a{
    hash[num] = true
   }
   result := []int{}
   for _,num := range b{
     if hash[num]{
        result = append(result, num)
     }
   }
   return result
}
```

### 7h. Difference

Write a function that returns elements in `a` that are NOT in `b`.

```
Example: difference([1,2,3,4], [3,4,5,6]) → [1,2]
```

```go
func difference(a, b []int) []int {
  hash := make(map[int]bool)
   for _,num := range b{
    hash[num] = true
   }
   result := []int{}
   for _,num := range a{
     if !hash[num]{
        result = append(result, num)
     }
   }
   return result

}
```

---

## Section 8: Stack Challenge

### 8a. Balanced Parentheses

Using your Stack, write a function that checks if parentheses are balanced.

```
"(())"    → true
"(()"     → false
"())("    → false
"()()()"  → true
""        → true
```

```go

type Stack struct {
	items []int
}

func NewStack() *Stack {
	return &Stack{
		items: []int{},
	}
}
func (c *Stack) Push(v int) {
	c.items = append(c.items, v)
}
func (c *Stack) Pop() (int, bool) {
	if len(c.items) == 0 {
		return 0, false
	}
	n := len(c.items)
	value := c.items[n-1]
	c.items = c.items[:n-1]
	return value, true
}
func (c *Stack) Peek() (int, bool) {
	if len(c.items) == 0 {
		return 0, false
	}
	return c.items[len(c.items)-1], true
}
func (c *Stack) IsEmpty() bool {
	return len(c.items) == 0
}
func (c *Stack) Size() int {
	return len(c.items)
}

func isBalanced(s string) bool {
  stack := NewStack()
  for _,ch := range s{
    if ch == "(". {
        stack.push(1)
    }else if ch == ")" {
        _, ok := stack.pop()
        if !ok {
            return false
        }
    }
  }
  return stack.IsEmpty()
}
```

### 8b. Reverse a String Using Stack

Write a function that reverses a string using a stack (push each character, then pop all).

```go
func reverseString(s string) string {
 stack := NewStack()

 for _,ch := range s{
    stack.push(int(ch))
 }
 str := ""
 for !stack.IsEmpty(){
    ch, _ := stack.pop()
    str += string(rune(ch))
 }
return str
}
```

---

## Section 9: True or False

Mark each statement as True or False, and explain WHY if false.

```
1. [ ] Slices are passed by value in Go. true

2. [ ] `append` always modifies the original slice. false, if cap exceeded new array is formed and orignal untouched

3. [ ] A nil slice and an empty slice (make([]int, 0)) behave the same with append. yes 

4. [ ] Sub-slices always share the same backing array as the original. yes 

5. [ ] `copy` returns the number of elements it attempted to copy. false, returns elements actually copied

6. [ ] Maps can be used as slice keys. No, Not camparable sot they cannot be

7. [ ] `len(s)` can never be greater than `cap(s)`. yes 

8. [ ] After `s = s[:0]`, the data is deleted from memory. no, the len becomes 0 but the backing array still there
```

---

## Section 10: Spot the Bug 🐛

### 10a.

```go
func removeDups(s []int) []int {
    seen := make(map[int]bool)
    for i, v := range s {
        if seen[v] {
            s = append(s[:i], s[i+1:]...)
        }
        seen[v] = true
    }
    return s
}
```

```
Bug:  append(s[:i], s[i+1:]...)
Fix: 
func removeDups(s []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    for i, v := range s {
        if !seen[v] {
            result = append(result,v)
            seen[v] = true
        }
      
    }
    return result
}
```

### 10b.

```go
func getFirst3(s []int) []int {
    return s[:3]
}

// called with: getFirst3([]int{1, 2})
```

```
Bug: return s[:3] out of index
Fix: 
func getFirst3(s []int) []int {
    n := 3
    if len(s) < n {
       n = return s
    }
    return s[:n]
}
```

### 10c.

```go
func double(s []int) {
    for i := range s {
        s[i] *= 2
    }
}

nums := []int{1, 2, 3}
double(nums)
fmt.Println(nums)
```

```
Output:[2,4,6]
Is this a bug? Why or why not: No bug as we are simply multpy by 2
```

---

_Great work! Slices are the backbone of Go programming. Master these and everything else becomes easier. 🚀_
