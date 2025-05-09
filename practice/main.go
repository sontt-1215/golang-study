package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type School struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Class struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	SchoolID          int    `json:"school_id"`
	HomeroomTeacherID int    `json:"homeroom_teacher_id"`
}

type Teacher struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Student struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type StudentClass struct {
	StudentID int `json:"student_id"`
	ClassID   int `json:"class_id"`
}

type Data struct {
	Schools        []School       `json:"schools"`
	Classes        []Class        `json:"classes"`
	Teachers       []Teacher      `json:"teachers"`
	Students       []Student      `json:"students"`
	StudentClasses []StudentClass `json:"student_classes"`
}

var reader *bufio.Reader

func loadData() (Data, error) {
	var data Data
	file, err := os.Open("data.json")
	if err != nil {
		return data, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return data, err
}

func readIntInput(prompt string) int {
	for {
		fmt.Println(prompt)
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Lỗi khi đọc input:", err)
			continue
		}
		text = strings.TrimSpace(text)
		num, err := strconv.Atoi(text)
		if err == nil {
			return num
		}
		fmt.Println("Vui lòng nhập số hợp lệ.")
	}
}

func getSortedClasses(classes []Class, teachers []Teacher, ch chan string) {
	sort.Slice(classes, func(i, j int) bool {
		return classes[i].Name < classes[j].Name
	})
	for _, c := range classes {
		teacherName := "N/A"
		for _, t := range teachers {
			if t.ID == c.HomeroomTeacherID {
				teacherName = t.Name
				break
			}
		}
		ch <- fmt.Sprintf("Class ID: %d, Tên: %s, GVCN: %s", c.ID, c.Name, teacherName)
	}
	close(ch)
}

func getSortedStudents(students []Student, ch chan string) {
	sort.Slice(students, func(i, j int) bool {
		return students[i].Name > students[j].Name
	})
	for _, s := range students {
		ch <- fmt.Sprintf("ID: %d, Tên: %s, Address: %s", s.ID, s.Name, s.Address)
	}
	close(ch)
}

func main() {
	reader = bufio.NewReader(os.Stdin)

	data, err := loadData()
	if err != nil {
		fmt.Println("Lỗi đọc file:", err)
		return
	}

	for {
		fmt.Println("\n=== MENU ===")
		fmt.Println("1. Danh sách lớp (theo tên ASC)")
		fmt.Println("2. Danh sách học sinh (theo tên DESC)")
		fmt.Println("3. Danh sách lớp học sinh đã tham gia")
		fmt.Println("4. Danh sách giáo viên với filter")
		fmt.Println("0. Thoát")

		choice := readIntInput("Chọn option:")

		switch choice {
		case 1:
			ch := make(chan string)
			go getSortedClasses(data.Classes, data.Teachers, ch)
			for line := range ch {
				fmt.Println(line)
			}

		case 2:
			ch := make(chan string)
			go getSortedStudents(data.Students, ch)
			for line := range ch {
				fmt.Println(line)
			}

		case 3:
			sid := readIntInput("Nhập student_id:")
			for _, sc := range data.StudentClasses {
				if sc.StudentID == sid {
					for _, c := range data.Classes {
						if c.ID == sc.ClassID {
							fmt.Printf("Class ID: %d, Tên: %s\n", c.ID, c.Name)
						}
					}
				}
			}

		case 4:
			fmt.Println("Filter:")
			fmt.Println("1. All")
			fmt.Println("2. Chủ nhiệm")
			fmt.Println("3. Có trên X học sinh")
			fmt.Println("4. Chủ nhiệm trên X lớp")
			f := readIntInput("Chọn filter:")

			switch f {
			case 1:
				for _, t := range data.Teachers {
					fmt.Printf("ID: %d, Tên: %s\n", t.ID, t.Name)
				}
			case 2:
				for _, t := range data.Teachers {
					for _, c := range data.Classes {
						if c.HomeroomTeacherID == t.ID {
							fmt.Printf("ID: %d, Tên: %s (GVCN)\n", t.ID, t.Name)
							break
						}
					}
				}
			case 3:
				x := readIntInput("Nhập X:")
				countMap := make(map[int]int)
				for _, sc := range data.StudentClasses {
					for _, c := range data.Classes {
						if c.ID == sc.ClassID {
							countMap[c.HomeroomTeacherID]++
						}
					}
				}
				for _, t := range data.Teachers {
					if countMap[t.ID] > x {
						fmt.Printf("ID: %d, Tên: %s (có %d học sinh)\n", t.ID, t.Name, countMap[t.ID])
					}
				}
			case 4:
				x := readIntInput("Nhập X:")
				countMap := make(map[int]int)
				for _, c := range data.Classes {
					countMap[c.HomeroomTeacherID]++
				}
				for _, t := range data.Teachers {
					if countMap[t.ID] > x {
						fmt.Printf("ID: %d, Tên: %s (chủ nhiệm %d lớp)\n", t.ID, t.Name, countMap[t.ID])
					}
				}
			}

		case 0:
			fmt.Println("Thoát chương trình.")
			return

		default:
			fmt.Println("Lựa chọn không hợp lệ.")
		}
	}
}

// docker compose build
// docker compose run --rm app
