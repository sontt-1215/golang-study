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
		fmt.Println(prompt)  // đảm bảo in ra trong container
		fmt.Print("> ")      // gọn gàng
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
			sort.Slice(data.Classes, func(i, j int) bool {
				return data.Classes[i].Name < data.Classes[j].Name
			})
			for _, c := range data.Classes {
				teacherName := "N/A"
				for _, t := range data.Teachers {
					if t.ID == c.HomeroomTeacherID {
						teacherName = t.Name
						break
					}
				}
				fmt.Printf("Class ID: %d, Name: %s, GVCN: %s\n", c.ID, c.Name, teacherName)
			}

		case 2:
			sort.Slice(data.Students, func(i, j int) bool {
				return data.Students[i].Name > data.Students[j].Name
			})
			for _, s := range data.Students {
				fmt.Printf("ID: %d, Name: %s, Address: %s\n", s.ID, s.Name, s.Address)
			}

		case 3:
			sid := readIntInput("Nhập student_id:")
			for _, sc := range data.StudentClasses {
				if sc.StudentID == sid {
					for _, c := range data.Classes {
						if c.ID == sc.ClassID {
							fmt.Printf("Class ID: %d, Name: %s\n", c.ID, c.Name)
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
					fmt.Printf("ID: %d, Name: %s\n", t.ID, t.Name)
				}
			case 2:
				for _, t := range data.Teachers {
					for _, c := range data.Classes {
						if c.HomeroomTeacherID == t.ID {
							fmt.Printf("ID: %d, Name: %s (GVCN)\n", t.ID, t.Name)
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
						fmt.Printf("ID: %d, Name: %s (có %d học sinh)\n", t.ID, t.Name, countMap[t.ID])
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
						fmt.Printf("ID: %d, Name: %s (chủ nhiệm %d lớp)\n", t.ID, t.Name, countMap[t.ID])
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
