package seeds

import (
	"log"
	"crm-go/config"
	"crm-go/models"
	"github.com/google/uuid"
)

func SeedTopics() {
	db := config.GetDB()


	// ✅ Validate UUIDs for GoLang Lessons
	topicID1, err01 := uuid.Parse("891b6645-db95-40dd-afff-8f68d8f34308")
	topicID2, err02 := uuid.Parse("a021a2fb-4df3-4e39-807d-fc42085a4b6d")
	topicID3, err03 := uuid.Parse("adb33026-ca07-4e25-8267-32a1020ccd90")
	topicID4, err04 := uuid.Parse("90325dd4-bc1e-4bf6-968b-9b83431de290")
	topicID5, err05 := uuid.Parse("fe1a5ac4-6111-4aeb-a1a4-b02650034a56")

	topicID6, err06 := uuid.Parse("991382ee-4a6c-44aa-a03a-2b4829e2acfe")
	topicID7, err07 := uuid.Parse("c96a12d1-3fd0-44c0-983c-f5442447f476")
	topicID8, err08 := uuid.Parse("a32f44d2-c319-4ff0-b391-cc0b6706ff93")
	topicID9, err09 := uuid.Parse("055db742-eef8-4f68-af62-47b3599bdd51")
	topicID10, err010 := uuid.Parse("52a5ca4d-bf6f-4ad4-acd9-7dcd5c1fa495")

	topicID11, err011 := uuid.Parse("1ab597d3-9867-4279-b2d7-ba093e36f090")
	topicID12, err012 := uuid.Parse("59861855-3bd3-4fd9-bc99-4210f4c0c5a0")
	topicID13, err013 := uuid.Parse("91bdc63f-79ce-42f0-bed5-14d95b0d7d36")
	topicID14, err014 := uuid.Parse("95ebe91c-f9c7-46c5-b58b-a98f295f6bc1")
	topicID15, err015 := uuid.Parse("8fde828b-9b7f-4c5d-8d81-811d7222195d")

	topicID16, err016 := uuid.Parse("21aa3941-a16b-4fa1-9c22-95e3e977e5f5")
	topicID17, err017 := uuid.Parse("252709a7-4f10-4147-90cb-a532657ecfa9")
	topicID18, err018 := uuid.Parse("9d9167c8-c408-4052-a0cf-8bb77c5b63c0")
	topicID19, err019 := uuid.Parse("d276a99d-95f0-4986-b084-f77ff2cb5fad")
	topicID20, err020 := uuid.Parse("414c6d5e-afb3-4e29-91fb-a840f281b477")

	topicID21, err021 := uuid.Parse("c1d73505-bedf-4470-8549-8808abefd902")
	topicID22, err022 := uuid.Parse("3a25d579-e144-4d21-9bcb-75b6ca9484aa")
	topicID23, err023 := uuid.Parse("93d8cc90-e126-4fd1-b443-7e51406e8a9e")
	topicID24, err024 := uuid.Parse("31ed53dd-c99f-41a5-ae92-ad23f1b4e3bc")
	topicID25, err025 := uuid.Parse("39428970-c68f-445b-9e26-4cd08365e834")

	if err01 != nil {
		log.Fatalf("❌ Invalid topic 1 UUID: %v", err01)
	}
	if err02 != nil {
		log.Fatalf("❌ Invalid topic 2 UUID: %v", err02)
	}
	if err03 != nil {
		log.Fatalf("❌ Invalid topic 3 UUID: %v", err03)
	}
	if err04 != nil {
		log.Fatalf("❌ Invalid topic 4 UUID: %v", err04)
	}
	if err05 != nil {
		log.Fatalf("❌ Invalid topic 5 UUID: %v", err05)
	}
	if err06 != nil {
		log.Fatalf("❌ Invalid topic 6 UUID: %v", err06)
	}
	if err07 != nil {
		log.Fatalf("❌ Invalid topic 7 UUID: %v", err07)
	}
	if err08 != nil {
		log.Fatalf("❌ Invalid topic 8 UUID: %v", err08)
	}
	if err09 != nil {
		log.Fatalf("❌ Invalid topic 9 UUID: %v", err09)
	}
	if err010 != nil {
		log.Fatalf("❌ Invalid topic 10 UUID: %v", err010)
	}
	if err011 != nil {
		log.Fatalf("❌ Invalid topic 11 UUID: %v", err011)
	}
	if err012 != nil {
		log.Fatalf("❌ Invalid topic 12 UUID: %v", err012)
	}
	if err013 != nil {
		log.Fatalf("❌ Invalid topic 13 UUID: %v", err013)
	}
	if err014 != nil {
		log.Fatalf("❌ Invalid topic 14 UUID: %v", err014)
	}
	if err015 != nil {
		log.Fatalf("❌ Invalid topic 15 UUID: %v", err015)
	}
	if err016 != nil {
		log.Fatalf("❌ Invalid topic 16 UUID: %v", err016)
	}
	if err017 != nil {
		log.Fatalf("❌ Invalid topic 17 UUID: %v", err017)
	}
	if err018 != nil {
		log.Fatalf("❌ Invalid topic 18 UUID: %v", err018)
	}
	if err019 != nil {
		log.Fatalf("❌ Invalid topic 19 UUID: %v", err019)
	}
	if err020 != nil {
		log.Fatalf("❌ Invalid topic 20 UUID: %v", err020)
	}
	if err021 != nil {
		log.Fatalf("❌ Invalid topic 21 UUID: %v", err021)
	}
	if err022 != nil {
		log.Fatalf("❌ Invalid topic 22 UUID: %v", err022)
	}
	if err023 != nil {
		log.Fatalf("❌ Invalid topic 23 UUID: %v", err023)
	}
	if err024 != nil {
		log.Fatalf("❌ Invalid topic 24 UUID: %v", err024)
	}
	if err025 != nil {
		log.Fatalf("❌ Invalid topic 25 UUID: %v", err025)
	}


	// ✅ Validate UUIDs for GoLang Lessons
	lessonID1, err1 := uuid.Parse("05d294c9-80f1-4203-8e78-d5f2c6d74c2c")
	lessonID2, err2 := uuid.Parse("00afe9fe-f686-4c04-8589-07e76beaebe4")
	lessonID3, err3 := uuid.Parse("5be66915-f2eb-49c2-b07a-41f467b1817a")
	lessonID4, err4 := uuid.Parse("813e6ca8-6513-4995-965f-3bb846f0b941")
	lessonID5, err5 := uuid.Parse("c6238448-c520-40ed-8a9b-54c1a5cd5c62")

	// ✅ Validate UUIDs for React Lessons
	lessonID6, err6 := uuid.Parse("5fbf5947-1904-40ef-b1c9-c1630b40e6b6")
	lessonID7, err7 := uuid.Parse("72223633-21be-4df6-9026-e4d7e09ba183")
	lessonID8, err8 := uuid.Parse("a778551d-fed1-47f1-a0b1-e4becd8cdc08")
	lessonID9, err9 := uuid.Parse("e9a14e98-761d-4d18-bd6a-db4e742413be")
	lessonID10, err10 := uuid.Parse("ff7edbb8-4515-49ba-b8b6-bac165bb61da")

	// ✅ Validate UUIDs for Python Lessons
	lessonID11, err11 := uuid.Parse("02770b59-bea7-4010-905f-02599cf25f33")
	lessonID12, err12 := uuid.Parse("47efa3d4-0b73-4c0a-8794-4542ba1847b6")
	lessonID13, err13 := uuid.Parse("648fcd8a-bc3a-4691-bd69-d365a73a463d")
	lessonID14, err14 := uuid.Parse("876e9404-448c-4b13-a3bf-9068ed888bcb")
	lessonID15, err15 := uuid.Parse("f3a4f79c-2a33-400a-b7f6-fea6f161d7f5")

	// ✅ Validate UUIDs for FastAPI Lessons
	lessonID16, err16 := uuid.Parse("9dbf26e9-0ae8-4a9e-90cc-178b25208f9d")
	lessonID17, err17 := uuid.Parse("d1c389a6-7396-45a4-a311-807e72b20d5c")
	lessonID18, err18 := uuid.Parse("da71896a-caf4-429f-b735-e7a847bb7069")
	lessonID19, err19 := uuid.Parse("deeefc60-f790-442d-afe8-51223f642411")
	lessonID20, err20 := uuid.Parse("fd84f5eb-c1cf-4d46-9b4b-118d43e294f8")

	// ✅ Validate UUIDs for Machine Learning Lessons
	lessonID21, err21 := uuid.Parse("49aff6ea-1ba0-410e-9563-a4053ee859db")
	lessonID22, err22 := uuid.Parse("a123f9a4-80ad-4336-ba52-beb76417a9c2")
	lessonID23, err23 := uuid.Parse("cd562cc9-37e1-4bc7-b7a5-8b711e411721")
	lessonID24, err24 := uuid.Parse("d9a7750e-42f4-4afc-8fb2-ded89ce0e2c8")
	lessonID25, err25 := uuid.Parse("ed8f20a1-0267-49d7-a8d4-6e232836423b")

	if err1 != nil {
		log.Fatalf("❌ Invalid lesson 1 UUID: %v", err1)
	}
	if err3 != nil {
		log.Fatalf("❌ Invalid lesson 3 UUID: %v", err3)
	}
	if err4 != nil {
		log.Fatalf("❌ Invalid lesson 4 UUID: %v", err4)
	}
	if err5 != nil {
		log.Fatalf("❌ Invalid lesson 5 UUID: %v", err5)
	}
	if err6 != nil {
		log.Fatalf("❌ Invalid lesson 6 UUID: %v", err6)
	}
	if err7 != nil {
		log.Fatalf("❌ Invalid lesson 7 UUID: %v", err7)
	}
	if err8 != nil {
		log.Fatalf("❌ Invalid lesson 8 UUID: %v", err8)
	}
	if err9 != nil {
		log.Fatalf("❌ Invalid lesson 9 UUID: %v", err9)
	}
	if err10 != nil {
		log.Fatalf("❌ Invalid lesson 10 UUID: %v", err10)
	}
	if err11 != nil {
		log.Fatalf("❌ Invalid lesson 11 UUID: %v", err11)
	}
	if err12 != nil {
		log.Fatalf("❌ Invalid lesson 12 UUID: %v", err12)
	}
	if err13 != nil {
		log.Fatalf("❌ Invalid lesson 13 UUID: %v", err13)
	}
	if err14 != nil {
		log.Fatalf("❌ Invalid lesson 14 UUID: %v", err14)
	}
	if err15 != nil {
		log.Fatalf("❌ Invalid lesson 15 UUID: %v", err15)
	}
	if err16 != nil {
		log.Fatalf("❌ Invalid lesson 16 UUID: %v", err16)
	}
	if err17 != nil {
		log.Fatalf("❌ Invalid lesson 17 UUID: %v", err17)
	}
	if err18 != nil {
		log.Fatalf("❌ Invalid lesson 18 UUID: %v", err18)
	}
	if err19 != nil {
		log.Fatalf("❌ Invalid lesson 19 UUID: %v", err19)
	}
	if err20 != nil {
		log.Fatalf("❌ Invalid lesson 20 UUID: %v", err20)
	}
	if err21 != nil {
		log.Fatalf("❌ Invalid lesson 21 UUID: %v", err21)
	}
	if err22 != nil {
		log.Fatalf("❌ Invalid lesson 22 UUID: %v", err22)
	}
	if err23 != nil {
		log.Fatalf("❌ Invalid lesson 23 UUID: %v", err23)
	}
	if err24 != nil {
		log.Fatalf("❌ Invalid lesson 24 UUID: %v", err24)
	}
	if err25 != nil {
		log.Fatalf("❌ Invalid lesson 25 UUID: %v", err25)
	}

		// ✅ Validate UUIDs  Courses
	courseId1, err26 := uuid.Parse("27d8ae14-4311-4380-8397-057ad5043fd6") // Go Lang Course
	courseId2, err27 := uuid.Parse("909b6026-30da-41f7-868f-42e6acba72c3") // React Course
	courseId3, err28 := uuid.Parse("b8ef3c14-d8ef-46fd-b63e-01b50cc9f227") // Python Course
	courseId4, err29 := uuid.Parse("c40c9c00-0779-490a-931b-e8dbd91549bf") // FastAPI Course
	courseId5, err30 := uuid.Parse("cdbe9a63-c659-4912-abb0-58dcb9d2f341") // Data Analysis Course

	if err26 != nil {
		log.Fatalf("❌ Invalid Course 26 UUID: %v", err26)
	}
	if err27 != nil {
		log.Fatalf("❌ Invalid Course 27 UUID: %v", err27)
	}
	if err28 != nil {
		log.Fatalf("❌ Invalid Course 28 UUID: %v", err28)
	}
	if err29 != nil {
		log.Fatalf("❌ Invalid Course 29 UUID: %v", err29)
	}
	if err30 != nil {
		log.Fatalf("❌ Invalid Course 30 UUID: %v", err30)
	}

	// ✅ Validate UUIDs for Go Lang Course
	moduleID31, err31 := uuid.Parse("1fb7c5b1-a438-4a19-9eb7-19a9046d0124")
	moduleID32, err32 := uuid.Parse("212021f5-c570-4898-b93d-129e4478497a")
	moduleID33, err33 := uuid.Parse("77af55f5-247c-44c7-b34d-67d118b95a26")
	moduleID34, err34 := uuid.Parse("b2df713a-0ca7-4fb3-9d43-d109b09d1491")
	moduleID35, err35 := uuid.Parse("ec980f1d-2ade-4e92-bb9b-912fd20c8f47")

	// ✅ Validate UUIDs for React
	moduleID36, err36 := uuid.Parse("1cee4f93-df85-4fb0-8557-c9773872e66a")
	moduleID37, err37 := uuid.Parse("79879e71-61fd-407d-8ccf-c3b4f813f9a0")
	moduleID38, err38 := uuid.Parse("7a4c17e4-3c69-433c-ad2d-f19e4ef01df3")
	moduleID39, err39 := uuid.Parse("bd26f9a3-4574-4fdb-9491-fb2be51927fb")
	moduleID40, err40 := uuid.Parse("bfe47ad2-49eb-4ed3-910f-3d22cb0e3e23")

	// ✅ Validate UUIDs for Python
	moduleID41, err41 := uuid.Parse("04936943-eaac-473c-a64d-3ac554465ea3")
	moduleID42, err42 := uuid.Parse("0beb9581-05ca-48f5-835c-4df7abefbb8f")
	moduleID43, err43 := uuid.Parse("353a4e83-e5b9-48f9-a7f7-d7ba351d6185")
	moduleID44, err44 := uuid.Parse("71c88341-9d44-4959-a272-5ea42a5f1503")
	moduleID45, err45 := uuid.Parse("96a5d5b4-abf0-46e3-b0c4-b35629966771")

	// ✅ Validate UUIDs for FastAPI
	moduleID46, err46 := uuid.Parse("2a2b4068-b6c5-4272-96c2-b0a76c7970eb")
	moduleID47, err47 := uuid.Parse("c94b72af-e65b-4203-91b8-c0fcdfaa44d5")
	moduleID48, err48 := uuid.Parse("dcd8dfc5-3e3b-4274-b442-6f0b1455268b")
	moduleID49, err49 := uuid.Parse("f2d65324-0e9d-42fe-8b50-26b67c1134d5")
	moduleID50, err50 := uuid.Parse("fe0d1f0d-3c07-4c96-a60a-80102aa1b6de")
	// ✅ Validate UUIDs for Machine Learning
	moduleID51, err51 := uuid.Parse("1f64e54d-0329-4cb3-8e4a-4199d7131723")
	moduleID52, err52 := uuid.Parse("490585ae-2c27-4295-be56-9ff5014e716c")
	moduleID53, err53 := uuid.Parse("89ab26e4-ff4d-484b-afcc-7c0a8c77ff7a")
	moduleID54, err54 := uuid.Parse("b9228452-d56a-4c98-9252-55280459326c")
	moduleID55, err55 := uuid.Parse("fb2aca67-4b51-40a2-897f-030b1b2b8af8")
	
	if err31 != nil {
		log.Fatalf("❌ Invalid module 31 UUID: %v", err31)
	}
	if err32 != nil {
		log.Fatalf("❌ Invalid module 32 UUID: %v", err32)
	}
	if err33 != nil {
		log.Fatalf("❌ Invalid module 33 UUID: %v", err33)
	}
	if err34 != nil {
		log.Fatalf("❌ Invalid module 34 UUID: %v", err34)
	}
	if err35 != nil {
		log.Fatalf("❌ Invalid module 35 UUID: %v", err35)
	}
	if err36 != nil {
		log.Fatalf("❌ Invalid module 36 UUID: %v", err36)
	}
	if err37 != nil {
		log.Fatalf("❌ Invalid module 37 UUID: %v", err37)
	}
	if err38 != nil {
		log.Fatalf("❌ Invalid module 38 UUID: %v", err38)
	}
	if err39 != nil {
		log.Fatalf("❌ Invalid module 39 UUID: %v", err39)
	}
	if err40 != nil {
		log.Fatalf("❌ Invalid module 40 UUID: %v", err40)
	}
	if err41 != nil {
		log.Fatalf("❌ Invalid module 41 UUID: %v", err41)
	}
	if err42 != nil {
		log.Fatalf("❌ Invalid module 42 UUID: %v", err42)
	}
	if err43 != nil {
		log.Fatalf("❌ Invalid module 43 UUID: %v", err43)
	}
	if err44 != nil {
		log.Fatalf("❌ Invalid module 44 UUID: %v", err44)
	}
	if err45 != nil {
		log.Fatalf("❌ Invalid module 45 UUID: %v", err45)
	}
	if err46 != nil {
		log.Fatalf("❌ Invalid module 46 UUID: %v", err46)
	}
	if err47 != nil {
		log.Fatalf("❌ Invalid module 47 UUID: %v", err47)
	}
	if err48 != nil {
		log.Fatalf("❌ Invalid module 48 UUID: %v", err48)
	}
	if err49 != nil {
		log.Fatalf("❌ Invalid module 49 UUID: %v", err49)
	}
	if err50 != nil {
		log.Fatalf("❌ Invalid module 50 UUID: %v", err50)
	}
	if err51 != nil {
		log.Fatalf("❌ Invalid module 51 UUID: %v", err51)
	}
	if err52 != nil {
		log.Fatalf("❌ Invalid module 52 UUID: %v", err52)
	}
	if err53 != nil {
		log.Fatalf("❌ Invalid module 53 UUID: %v", err53)
	}
	if err54 != nil {
		log.Fatalf("❌ Invalid module 54 UUID: %v", err54)
	}
	if err55 != nil {
		log.Fatalf("❌ Invalid module 55 UUID: %v", err55)	
	}



	createdByID, err2 := uuid.Parse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")

	if err2 != nil {
		log.Fatalf("❌ Invalid CreatedBy UUID: %v", err2)
	}



	lessons := []models.Topics{

		// Go Lang Course
		{
			ID:        topicID1,
			LessonID:        lessonID1,
			CourseID: courseId1,
			ModuleID:   moduleID31,
			TutorID:  createdByID,
			Title:     "Overview of the Go Language",
			ContentType:     "article",
			ContentURL: "https://example.com/go-introduction",
			ContentText: "Go, also known as Golang, is an open-source programming language designed at Google to build simple, reliable, and efficient software. It is widely used for backend development, cloud services, and distributed systems. Go's syntax is clean and concise, making it easy to learn for developers coming from other languages. It features built-in support for concurrent programming, allowing developers to write programs that can efficiently utilize multiple CPU cores. Go also has a strong standard library that provides a wide range of functionalities, from web servers to cryptography. With its focus on performance and simplicity, Go has become a popular choice for building scalable and high-performance applications.",
			Order:      1,
		},
		{
			ID:        topicID2,
			LessonID:        lessonID2,
			CourseID: courseId1,
			ModuleID:   moduleID32,
			TutorID:  createdByID,
			Title:     "Installing and Verifying Go",
			ContentType:     "video",
			ContentURL: "https://example.com/go-installation-video",
			ContentText: "To install Go, download the installer from the official Go website and follow the installation instructions for your operating system. After installation, run go version in the terminal to confirm that Go has been installed successfully. You should see the installed Go version displayed, indicating that Go is ready to use on your system.",
			Order:      2,
		},
			{
			ID:        topicID3,
			LessonID:        lessonID3,
			CourseID: courseId1,
			ModuleID:   moduleID33,
			TutorID:  createdByID,
			Title:     "Declaring Variables in Go",
			ContentType:     "article",
			ContentURL: "https://example.com/go-installation-video",
			ContentText: "Variables in Go can be declared using the var keyword or the shorthand := syntax. Go also supports several basic data types such as integers, floats, strings, and booleans. Example: <> var name string = 'Gopher' age := 10 </>" ,
			Order:      3,
		},
		{
			ID: 		topicID4,
			LessonID:        lessonID4,
			CourseID: courseId1,
			ModuleID:   moduleID34,
			TutorID:  createdByID,
			Title:     "Using If Statements in Go",
			ContentType:     "article",
			ContentURL: "https://example.com/go-if-statements",
			ContentText: "The if statement in Go allows a program to execute code based on condition Example age :=  if age >= 18 { fmt.Println('You are an adult')}" ,
			Order:      4,
		},
		{
			ID: 		topicID5,
			LessonID:        lessonID5,
			CourseID: courseId1,
			ModuleID:   moduleID35,
			TutorID:  createdByID,
			Title:     "Writing Your First Function in Go",
			ContentType:     "article",
			ContentURL: "https://example.com/go-functions",
			ContentText: "Functions in Go allow developers to group related code together. A function is defined using the func keyword. Example: func greet(name string) string { return 'Hello' + name} Calling the function: message := greet('Emmanuel') fmt.Println(message)" ,
			Order:      5,
		},

			// React Course
		{
			ID:        topicID6,
			LessonID:        lessonID6,
			CourseID: courseId2,
			ModuleID:   moduleID36,
			TutorID:  createdByID,
			Title:     "Overview of React",
			ContentType:     "article",
			ContentURL: "https://example.com/react-introduction",
			ContentText: "React is a JavaScript library used for building user interfaces, especially single-page applications. It allows developers to build reusable UI components and efficiently update the user interface using a virtual DOM. React is maintained by Facebook and has a large community of developers. It uses a component-based architecture, where each component manages its own state and can be composed to create complex UIs. React also supports features like hooks, which allow developers to use state and other React features without writing class components. With its focus on performance and developer experience, React has become one of the most popular libraries for frontend development." ,
			Order:      1,
		},
		{
			ID:        topicID7,
			LessonID:        lessonID7,
			CourseID: courseId2,
			ModuleID:   moduleID37,
			TutorID:  createdByID,
			Title:     "Creating a React App with Vite",
			ContentType:     "article",
			ContentURL: "https://example.com/react-vite",
			ContentText: "To create a React project using Vite, run the following command in your terminal: npm create vite@latest my-react-app Then install dependencies and start the development server. cd my-react-ap npm instal npm run dev This will start the React development server and open the project in your browser. Vite is a build tool that provides a faster and leaner development experience for modern web projects. It offers instant server start, fast hot module replacement (HMR), and optimized production builds. When creating a React app with Vite, developers can leverage its features to enhance their development workflow and improve the overall performance of their applications." ,
			Order:      2,
		},
			{
			ID:        topicID8,
			LessonID:        lessonID8,
			CourseID: courseId2,
			ModuleID:   moduleID38,
			TutorID:  createdByID,
			Title:     "Creating Components",
			ContentType:     "article",
			ContentURL: "https://example.com/react-component",
			ContentText: "Components are the building blocks of a React application. A functional component is simply a JavaScript function that returns JSX. Example: function Welcome() { return <h1>Hello, React!</h1>; } export default Welcome; Components help organize UI and make code reusable." ,
			Order:      3,
		},
		{
			ID:        topicID9,
			LessonID:        lessonID9,
			CourseID: courseId2,
			ModuleID:   moduleID39,
			TutorID:  createdByID,
			Title:     "Passing Data with Props",
			ContentType:     "article",
			ContentURL: "https://example.com/react-props",
			ContentText: "Props allow components to receive data from their parent components.Example: function Greeting(props) { return <h1>Hello {props.name}</h1>;} <Greeting name='Emmanuel' /> Props make components flexible and reusable. Props (short for properties) are a way to pass data from a parent component to a child component in React. They are read-only and help in making components reusable and maintainable." ,
			Order:      4,
		},
		{
			ID:        topicID10,
			LessonID:        lessonID10,
			CourseID: courseId2,
			ModuleID:   moduleID40,
			TutorID:  createdByID,
			Title:     "Handling Click Events",
			ContentType:     "article",
			ContentURL: "https://example.com/react-react-events",
			ContentText: "React handles events similar to HTML but uses camelCase syntax. Example: function Button() { function handleClick() {alert('Button clicked!');} return <button onClick={handleClick}>Click Me</button>;} Event handling allows React applications to respond to user actions. Event handling in React is similar to handling events in regular HTML, but with some differences. In React, event handlers are functions that are called when an event occurs. For example, to handle a button click, you would define a function and pass it to the button's onClick prop." ,
			Order:      5,
		},

		// Python Course
		{
			ID:        topicID11,
			LessonID:        lessonID11,
			CourseID: courseId3,
			ModuleID:   moduleID41,
			TutorID:  createdByID,
			Title:     "Overview of Python",
			ContentType:     "article",
			ContentURL: "https://example.com/python-introduction",
			ContentText: "Python is a high-level programming language known for its simple syntax and readability. It is widely used for web development, automation, machine learning, and data analysis. Python allows developers to write powerful programs with fewer lines of code compared to many other programming languages." ,
			Order:      1,
		},
		{
			ID:        topicID12,
			LessonID:        lessonID12,
			CourseID: courseId3,
			ModuleID:   moduleID42,
			TutorID:  createdByID,
			Title:     "Installing and Running Python",
			ContentType:     "article",
			ContentURL: "https://example.com/python-installation",
			ContentText: "To install Python, visit the official Python website and download the latest version for your operating system. Once installed, you can run Python scripts from the command line or use an integrated development environment (IDE) like PyCharm or Visual Studio Code. To install Python, download the installer from the official Python website and follow the installation steps for your operating system. After installation, open the terminal and run: python --version This command confirms that Python has been successfully installed on your system." ,
			Order:      2,
		},
			{
			ID:        topicID13,
			LessonID:        lessonID13,
			CourseID: courseId3,
			ModuleID:   moduleID43,
			TutorID:  createdByID,
			Title:     "Declaring Variables in Python",
			ContentType:     "article",
			ContentURL: "https://example.com/python-variables",
			ContentText: "In Python, variables are declared by simply assigning a value to a name. For example, x = 5 declares a variable named x and assigns it the value 5. Variables in Python are dynamically typed, meaning you don't need to explicitly declare their type." ,
			Order:      3,
		},
		{
			ID:        topicID14,
			LessonID:        lessonID14,
			CourseID: courseId3,
			ModuleID:   moduleID44,
			TutorID:  createdByID,
			Title:     "Using If Statements in Python",
			ContentType:     "article",
			ContentURL: "https://example.com/python-if-statements",
			ContentText: "In Python, if statements are used to execute code based on certain conditions. The basic syntax is: if condition: # code to execute if condition is true. Conditional statements allow Python programs to make decisions. Example: age = 18 if age >= 18: print('You are an adult') else: print('You are a minor') These conditions help programs respond to different inputs and situations." ,
			Order:      4,
		},
		{
			ID:        topicID15,
			LessonID:        lessonID15,
			CourseID: courseId3,
			ModuleID:   moduleID45,
			TutorID:  createdByID,
			Title:     "Writing Your First Function in Python",
			ContentType:     "article",
			ContentURL: "https://example.com/python-functions",
			ContentText: "In Python, functions are defined using the def keyword. Example: def greet(name): print(f'Hello, {name}!') This function takes a parameter and prints a greeting. Functions help organize code into reusable blocks." ,
			Order:      5,
		},

			// Fast API Course
		{
			ID:        topicID16,
			LessonID:        lessonID16,
			CourseID: courseId4,
			ModuleID:   moduleID46,
			TutorID:  createdByID,
			Title:     "Overview of FastAPI",
			ContentType:     "article",
			ContentURL: "https://example.com/fastapi-introduction",
			ContentText: "FastAPI is a modern web framework for building APIs with Python. It is known for its speed, simplicity, and automatic API documentation. FastAPI uses Python type hints to validate data and generate interactive API documentation automatically." ,
			Order:      1,
		},
		{
			ID:        topicID17,
			LessonID:        lessonID17,
			CourseID: courseId4,
			ModuleID:   moduleID47,
			TutorID:  createdByID,
			Title:     "Installing FastAPI and Uvicorn",
			ContentType:     "video",
			ContentURL: "https://example.com/install-fastapi-video",
			ContentText: "To install FastAPI, you can use pip, the Python package installer. Run the following command in your terminal: pip install fastapi[all] This command will install FastAPI along with all its dependencies, including Uvicorn, which is an ASGI server used to run FastAPI applications. To install FastAPI, use pip to install the required packages: pip install fastapi uvicorn After installation, you can start building APIs and run them using the Uvicorn server." ,
			Order:      2,
		},
			{
			ID:        topicID18,
			LessonID:        lessonID18,
			CourseID: courseId4,
			ModuleID:   moduleID48,
			TutorID:  createdByID,
			Title:     "Creating a Simple GET Endpoint",
			ContentType:     "article",
			ContentURL: "https://example.com/fastapi-get-endpoint",
			ContentText: "In FastAPI, you can create API endpoints using the @app.get(), @app.post(), and other decorators. Example: @app.get('/users') def get_users(): return {'users': []} This endpoint handles GET requests to the /users route." ,
			Order:      3,
		},
		{
			ID:        topicID19,
			LessonID:        lessonID19,
			CourseID: courseId4,
			ModuleID:   moduleID49,
			TutorID:  createdByID,
			Title:     "Path Parameters in FastAPI",
			ContentType:     "article",
			ContentURL: "https://example.com/fastapi-path-parameters",
			ContentText: "In FastAPI, you can access request data using path parameters, query parameters, and request bodies. Example: @app.get('/users/{user_id}') def get_user(user_id: int): return {'user_id': user_id} This endpoint handles GET requests to the /users/{user_id} route." ,
			Order:      4,
		},
		{
			ID:        topicID20,
			LessonID:        lessonID20,
			CourseID: courseId4,
			ModuleID:   moduleID50,
			TutorID:  createdByID,
			Title:     "Defining a Pydantic Model",
			ContentType:     "article",
			ContentURL: "https://example.com/fastapi-pydantic-models",
			ContentText: "FastAPI uses Pydantic models to validate incoming request data automatically. Example: from pydantic import BaseModel class User(BaseModel): name: str age: int @app.post('/users') def create_user(user: User): return user This ensures that incoming data matches the expected structure before processing it." ,
			Order:      5,
		},
		
			// Data Science Course
		{
			ID:        topicID21,
			LessonID:        lessonID21,
			CourseID: courseId5,
			ModuleID:   moduleID51,
			TutorID:  createdByID,
			Title:     "Overview of Data Analysis",
			ContentType:     "article",
			ContentURL: "https://example.com/data-analysis-introduction",
			ContentText: "Data analysis is the process of inspecting, cleaning, transforming, and modeling data to discover useful information, draw conclusions, and support decision-making. It is widely used in business, science, and engineering." ,
			Order:      1,
		},
		{
			ID:        topicID22,
			LessonID:        lessonID22,
			CourseID: courseId5,
			ModuleID:   moduleID52,
			TutorID:  createdByID,
			Title:     "Handling Missing Data and Data Cleaning",
			ContentType:     "article",
			ContentURL: "https://example.com/data-cleaning",
			ContentText: "Before analyzing data, it is important to clean it. Missing or inconsistent values can lead to incorrect results. Techniques include removing missing values, filling them with defaults, or using statistical methods to impute missing data. Example in Python: import pandas as pd df = pd.read_csv('data.csv') df.fillna(0, inplace=True)",
			Order:      2,
		},
			{
			ID:        topicID23,
			LessonID:        lessonID23,
			CourseID: courseId5,
			ModuleID:   moduleID53,
			TutorID:  createdByID,
			Title:     "Descriptive Statistics and Visualization",
			ContentType:     "video",
			ContentURL: "https://example.com/eda-video",
			ContentText: "Exploratory Data Analysis (EDA) helps uncover patterns in data. Common techniques include calculating mean, median, standard deviation, and plotting histograms, scatter plots, and box plots using libraries like Pandas and Matplotlib. Example: import pandas as pd import matplotlib.pyplot as plt df = pd.read_csv('data.csv') df['age'].hist() plt.show()",
			Order:      3,
		},
		{
			ID:        topicID24,
			LessonID:        lessonID24,
			CourseID: courseId5,
			ModuleID:   moduleID54,
			TutorID:  createdByID,
			Title:     "Data Manipulation with Pandas and NumPy",
			ContentType:     "article",
			ContentURL: "https://example.com/pandas-data-manipulation",
			ContentText: "Pandas is a powerful library for data manipulation and analysis. You can read, filter, group, and summarize data easily. Example: import pandas as pddf = pd.read_csv('data.csv') print(df.head()) print(df['salary'].mean())",
			Order:      4,
		},
		{
			ID:        topicID25,
			LessonID:        lessonID25,
			CourseID: courseId5,
			ModuleID:   moduleID55,
			TutorID:  createdByID,
			Title:     "Visualizing Data with Matplotlib and Seaborn",
			ContentType:     "video",
			ContentURL: "https://example.com/matplotlib-visualization",
			ContentText: "Visualizations make it easier to interpret and communicate data insights. Matplotlib is a Python library that can create line plots, bar charts, scatter plots, and more. Example: import matplotlib.pyplot as plt x = [1, 2, 3, 4] y = [10, 20, 15, 25] plt.plot(x, y) plt.title('Sample Line Chart') plt.xlabel('X Axis') plt.ylabel('Y Axis') plt.show()",
			Order:      5,
		},
		
	
		
	}

	for _, lesson := range lessons {
		if err := db.Create(&lesson).Error; err != nil {
			log.Printf("❌ Failed to seed Lesson: %v", err)
		} else {
			log.Printf("✅ Seeded Lesson: %v", &lesson.Title)
		}
	}
}
