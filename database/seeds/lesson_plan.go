package seeds

import (
	"log"

	"github.com/google/uuid"
	"fmt"
	"crm-go/models"
	"crm-go/config"

)

func SeedLessonPlans() error{

	db := config.GetDB()

	lessonPlanID1 := uuid.MustParse("a6a10726-b1cd-4d62-b299-742f1f0cc01d")
	lessonPlanID2 := uuid.MustParse("edc841d4-dec9-48ae-a50e-c03de738d1c0")
	lessonPlanID3 := uuid.MustParse("6a0d2d35-9a3a-446d-b2f0-bfbea032dd7d")
	lessonPlanID4 := uuid.MustParse("b25fa7c4-346b-42eb-87be-8cfc8c9fd3e2")
	lessonPlanID5 := uuid.MustParse("459d5b24-2d26-456f-824c-3e646825770d")
	lessonPlanID6 := uuid.MustParse("04ce4738-33d1-4424-a097-87e952b4310f")
	lessonPlanID7 := uuid.MustParse("4e697776-d552-43a5-aeea-31dd5a43daf1")
	lessonPlanID8 := uuid.MustParse("d203d65f-efc1-4ca9-86e2-2eab351f99c1")
	lessonPlanID9 := uuid.MustParse("44ec601f-d248-469d-910e-e20f6b670c8c")
	lessonPlanID10 := uuid.MustParse("d0cfc75d-9cbf-45bb-835b-96ef5beee9af")
	lessonPlanID11 := uuid.MustParse("0ee5a640-c81c-44d0-ae5d-234c4d16ad78")
	lessonPlanID12 := uuid.MustParse("3873bf28-d1f4-4b16-b709-335bc3574efb")
	lessonPlanID13 := uuid.MustParse("879f633a-692c-434f-8e2b-eff57a480547")

	adminID := uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd",)

	lessonPlans := []models.LessonPlan{
		{
			ID:       lessonPlanID1,

			LessonID: uuid.MustParse("e2a83449-fd98-4b24-8fcf-e51b4e17456b"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to define micro-organisms, identify the major groups of micro-organisms, describe their characteristics, and state examples of useful and harmful micro-organisms.",

			Introduction: "The teacher asks students what they understand by germs and microscopic organisms. Students share their ideas before the teacher introduces the topic of micro-organisms.",

			PreviousKnowledge: "Students have basic knowledge of living organisms and understand that plants and animals are made up of cells.",

			LessonContent: "Meaning of micro-organisms; major groups including bacteria, fungi, protozoa, algae and viruses; characteristics of micro-organisms; useful and harmful effects of micro-organisms.",

			TeacherActivities: "The teacher introduces the concept of micro-organisms, explains their major groups, displays charts or images, gives examples, and guides students in identifying useful and harmful micro-organisms.",

			StudentActivities: "Students listen to the explanation, observe diagrams, identify different micro-organisms, answer questions, participate in discussions, and record important points.",

			TeachingAids: "Biology textbook, charts, diagrams, microscope images, whiteboard and marker.",

			Evaluation: "What are micro-organisms? Name four groups of micro-organisms. Give two examples of useful micro-organisms and two examples of harmful micro-organisms.",

			Conclusion: "The teacher summarizes the meaning, groups, characteristics and importance of micro-organisms and emphasizes that some are beneficial while others can cause disease.",

			CreatedBy: adminID,
		},

		{
			ID:       lessonPlanID2,

			LessonID: uuid.MustParse("9e12aaf1-fdb9-430e-8f5e-e98c36ec3bf9"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to explain how micro-organisms are cultured, describe basic methods of culturing micro-organisms, and identify micro-organisms using observable characteristics.",

			Introduction: "The teacher reviews the previous lesson by asking students to name different groups of micro-organisms and then introduces the idea that micro-organisms can be grown under suitable conditions.",

			PreviousKnowledge: "Students know the major groups of micro-organisms and understand that micro-organisms require suitable conditions for growth.",

			LessonContent: "Meaning of culture; conditions required for growth; basic culturing techniques; culture media; observation and identification of micro-organisms; laboratory safety.",

			TeacherActivities: "The teacher explains the concept of culturing, describes suitable growth conditions, demonstrates safe laboratory procedures using diagrams or prepared materials, and explains identification techniques.",

			StudentActivities: "Students observe the demonstration, identify conditions necessary for microbial growth, discuss observations, answer questions, and write notes.",

			TeachingAids: "Biology textbook, laboratory diagrams, prepared culture images, charts, whiteboard and marker.",

			Evaluation: "What is a microbial culture? State three conditions required for microbial growth. Mention two ways micro-organisms can be identified.",

			Conclusion: "The teacher reviews the process of culturing and identification and reminds students that proper laboratory safety must always be observed.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID3,

			LessonID: uuid.MustParse("d837b61a-d747-4cb9-bca8-c165bc81e6fe"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to define sexually transmitted infections, identify common STIs, describe their modes of transmission, and state ways of preventing transmission.",

			Introduction: "The teacher introduces the topic through a discussion about diseases that can be transmitted from one person to another through sexual contact.",

			PreviousKnowledge: "Students understand the basic concepts of infectious diseases and how diseases can be transmitted between individuals.",

			LessonContent: "Meaning of sexually transmitted infections; examples including HIV/AIDS, gonorrhoea, syphilis, chlamydia, genital herpes and hepatitis B; modes of transmission; signs and symptoms; prevention.",

			TeacherActivities: "The teacher explains STIs using age-appropriate language, discusses common examples, explains modes of transmission, and emphasizes prevention and responsible health practices.",

			StudentActivities: "Students listen, participate in discussions, identify examples of STIs, answer questions, and record prevention measures.",

			TeachingAids: "Biology textbook, health education charts, diagrams, whiteboard and marker.",

			Evaluation: "What are STIs? Name four sexually transmitted infections. State three ways through which STIs can be transmitted and three ways they can be prevented.",

			Conclusion: "The teacher emphasizes that STIs can have serious consequences and that prevention, accurate information and responsible health decisions are important.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID4,

			LessonID: uuid.MustParse("caafaae2-7efa-4cd8-89fb-a893c35a88e2"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to explain beneficial and harmful activities of micro-organisms and describe their roles in food production, decomposition, medicine and disease.",

			Introduction: "The teacher asks students whether all micro-organisms are harmful and allows students to discuss examples of useful micro-organisms.",

			PreviousKnowledge: "Students understand the major groups of micro-organisms and know that some micro-organisms cause disease.",

			LessonContent: "Useful activities of micro-organisms; decomposition; fermentation; food production; medicine; nitrogen cycle; harmful activities including food spoilage and disease.",

			TeacherActivities: "The teacher explains the activities of micro-organisms, provides everyday examples, uses diagrams to illustrate decomposition and fermentation, and discusses their beneficial and harmful effects.",

			StudentActivities: "Students provide examples from everyday life, participate in group discussions, classify activities as beneficial or harmful, and answer questions.",

			TeachingAids: "Textbook, charts, pictures of fermented foods, diagrams and whiteboard.",

			Evaluation: "State four useful activities of micro-organisms. Mention three harmful effects of micro-organisms. Explain the role of micro-organisms in decomposition.",

			Conclusion: "The teacher summarizes how micro-organisms affect human life, agriculture, food production and the environment.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID5,

			LessonID: uuid.MustParse("3f93af4b-9667-4772-8ada-b53875136304"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to identify major causes of STIs, explain their effects on individuals and society, and describe appropriate preventive measures.",

			Introduction: "The teacher reviews the meaning of STIs and asks students to recall some examples discussed in the previous lesson.",

			PreviousKnowledge: "Students understand the meaning, examples and basic modes of transmission of sexually transmitted infections.",

			LessonContent: "Causes and risk factors of STIs; effects on reproductive health; complications; social and economic effects; prevention; testing and treatment; importance of seeking medical advice.",

			TeacherActivities: "The teacher explains the causes and effects of STIs, discusses risk factors, presents preventive strategies, and corrects common misconceptions.",

			StudentActivities: "Students participate in discussions, identify risk factors, explain possible effects, answer questions, and list preventive measures.",

			TeachingAids: "Health education charts, textbook, diagrams, whiteboard and marker.",

			Evaluation: "State four causes or risk factors of STIs. Mention four effects of untreated STIs. State five measures for preventing STIs.",

			Conclusion: "The teacher reinforces the importance of prevention, early testing, appropriate treatment and responsible health practices.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID6,

			LessonID: uuid.MustParse("710f0501-9c6b-49f7-b8b9-5b1ef22f519f"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to explain methods of controlling harmful micro-organisms and distinguish between physical, chemical and biological control methods.",

			Introduction: "The teacher asks students how food, drinking water and medical equipment can be protected from harmful micro-organisms.",

			PreviousKnowledge: "Students know that some micro-organisms cause diseases and spoil food.",

			LessonContent: "Control of micro-organisms; sterilization; disinfection; antiseptics; heat treatment; refrigeration; drying; use of chemicals; hygiene and sanitation.",

			TeacherActivities: "The teacher explains different control methods, gives everyday examples, demonstrates proper hygiene practices, and compares sterilization and disinfection.",

			StudentActivities: "Students identify methods used to control micro-organisms at home, in hospitals and in food industries and participate in classroom discussions.",

			TeachingAids: "Textbook, charts, pictures, disinfectant labels, whiteboard and marker.",

			Evaluation: "What is sterilization? Differentiate between sterilization and disinfection. State five methods of controlling harmful micro-organisms.",

			Conclusion: "The teacher summarizes the major control methods and explains their importance in preventing disease and food spoilage.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID7,

			LessonID: uuid.MustParse("90169597-ef87-4451-8d79-e6f3451b344d"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to define plant classification, identify major plant groups, describe their characteristics, and give examples of plants in each group.",

			Introduction: "The teacher displays pictures of different plants and asks students to identify similarities and differences among them.",

			PreviousKnowledge: "Students know common plants found in their environment and understand basic plant structures.",

			LessonContent: "Meaning and importance of classification; major plant groups; algae; bryophytes; pteridophytes; gymnosperms; angiosperms; monocotyledons and dicotyledons.",

			TeacherActivities: "The teacher explains the basis of plant classification, displays plant specimens or diagrams, and guides students in grouping plants according to their characteristics.",

			StudentActivities: "Students observe specimens, identify characteristics, classify plants, participate in group work, and complete classification exercises.",

			TeachingAids: "Plant specimens, charts, textbook, diagrams, photographs, whiteboard and marker.",

			Evaluation: "Why are plants classified? Name five major plant groups. State two characteristics of flowering plants.",

			Conclusion: "The teacher summarizes the major plant groups and explains that classification helps scientists organize and study plant diversity.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID8,

			LessonID: uuid.MustParse("e6492333-05af-4d29-8e38-8f83d050304d"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to identify common plant pests and diseases, describe their effects on crops, and recognize common signs of plant disease.",

			Introduction: "The teacher asks students to describe problems they have observed on crops, such as damaged leaves, wilting or discoloration.",

			PreviousKnowledge: "Students understand basic plant structures and the importance of plants in agriculture.",

			LessonContent: "Meaning of plant pests and diseases; common insect pests; fungal, bacterial and viral diseases; symptoms; effects on crop production; examples of affected crops.",

			TeacherActivities: "The teacher displays pictures of affected plants, explains symptoms, identifies common pests and diseases, and discusses their effects on agricultural production.",

			StudentActivities: "Students observe pictures or specimens, identify symptoms, discuss causes, and classify examples of pests and diseases.",

			TeachingAids: "Affected plant specimens, charts, photographs, textbook, diagrams and whiteboard.",

			Evaluation: "What is a plant pest? Name four common plant pests. State three signs of plant disease. Mention three effects of plant diseases on crop production.",

			Conclusion: "The teacher summarizes common plant pests and diseases and explains why early identification is important for crop protection.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID9,

			LessonID: uuid.MustParse("ae02e07e-e0f9-4d58-aa29-1f7f1d35e382"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to identify common animal pests and diseases, describe their effects on livestock, and recognize common signs of animal disease.",

			Introduction: "The teacher asks students to name common livestock animals and discuss health problems they may have observed in animals.",

			PreviousKnowledge: "Students have basic knowledge of farm animals and understand that animals can suffer from diseases.",

			LessonContent: "Animal pests; internal and external parasites; common livestock diseases; bacterial, viral, fungal and parasitic diseases; signs and effects of animal diseases.",

			TeacherActivities: "The teacher explains common animal pests and diseases, displays relevant diagrams, discusses symptoms, and explains how diseases affect livestock production.",

			StudentActivities: "Students identify common pests and diseases, observe diagrams, discuss symptoms, and answer questions.",

			TeachingAids: "Agricultural Science textbook, animal health charts, photographs, diagrams and whiteboard.",

			Evaluation: "What is an animal pest? Give four examples of animal pests. Name four livestock diseases. State three signs of disease in farm animals.",

			Conclusion: "The teacher emphasizes the importance of recognizing pests and diseases early to reduce livestock losses.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID10,

			LessonID: uuid.MustParse("ebe29753-3c96-4a0b-a7c0-5539d4d63a0b"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to explain methods of controlling animal diseases, describe preventive measures, and explain the importance of good livestock management.",

			Introduction: "The teacher reviews common animal diseases and asks students how farmers can prevent their animals from becoming sick.",

			PreviousKnowledge: "Students know common animal pests and diseases and can identify some signs of disease.",

			LessonContent: "Prevention and control of animal diseases; vaccination; quarantine; sanitation; proper feeding; vector control; veterinary care; isolation of sick animals.",

			TeacherActivities: "The teacher explains disease prevention methods, discusses vaccination and quarantine, and describes good livestock management practices.",

			StudentActivities: "Students discuss preventive measures, identify appropriate control methods, answer questions, and complete classroom exercises.",

			TeachingAids: "Agricultural Science textbook, livestock health charts, diagrams and whiteboard.",

			Evaluation: "State five methods of controlling animal diseases. What is quarantine? Why is vaccination important? Mention three good livestock management practices.",

			Conclusion: "The teacher summarizes the main methods of controlling animal diseases and emphasizes prevention as an important part of livestock management.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID11,

			LessonID: uuid.MustParse("e4698004-cad7-4c38-979c-4d0e78101bc6"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to define food production, identify sources of food, explain factors affecting food production, and describe basic methods of producing major food crops and animal products.",

			Introduction: "The teacher asks students to list foods commonly consumed at home and identify where each food comes from.",

			PreviousKnowledge: "Students understand the importance of agriculture and know common crops and livestock raised by farmers.",

			LessonContent: "Meaning of food production; crop production; livestock production; fisheries; sources of food; factors affecting food production; importance of food production.",

			TeacherActivities: "The teacher explains food production, discusses crop and animal production systems, identifies factors affecting production, and uses local examples.",

			StudentActivities: "Students identify food sources, participate in discussions, classify foods according to their sources, and answer questions.",

			TeachingAids: "Agricultural Science textbook, food samples or pictures, charts, diagrams and whiteboard.",

			Evaluation: "What is food production? Identify four sources of food. State five factors that affect food production. Explain two methods of food production.",

			Conclusion: "The teacher summarizes the major sources and methods of food production and highlights the importance of agriculture to food security.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID12,

			LessonID: uuid.MustParse("3c8f106c-0203-4c69-bfb2-f507b1fc5819"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to explain food preservation, identify common preservation methods, explain the importance of proper food storage, and differentiate between preservation and storage.",

			Introduction: "The teacher asks students how their families keep food from spoiling when it cannot be consumed immediately.",

			PreviousKnowledge: "Students understand that food can spoil and that micro-organisms contribute to food spoilage.",

			LessonContent: "Meaning of food preservation and storage; drying; smoking; salting; refrigeration; freezing; canning; fermentation; pasteurization; proper storage conditions.",

			TeacherActivities: "The teacher explains different preservation methods, demonstrates examples, discusses the principles behind preservation, and explains proper food storage practices.",

			StudentActivities: "Students identify preservation methods used at home, classify foods according to preservation methods, participate in discussions, and answer questions.",

			TeachingAids: "Food samples, pictures, textbook, charts, refrigerator/freezer illustrations, whiteboard and marker.",

			Evaluation: "What is food preservation? State six methods of preserving food. Explain how drying preserves food. State four proper food storage practices.",

			Conclusion: "The teacher summarizes food preservation and storage methods and explains their importance in reducing food spoilage and waste.",

			CreatedBy: adminID,

		},

		{
			ID:       lessonPlanID13,

			LessonID: uuid.MustParse("9983316d-d951-4ba5-8959-90faa942dd49"),

			BehaviouralObjectives: "By the end of the lesson, students should be able to recall major concepts covered during the term, explain key biological and agricultural processes, answer revision questions, and identify areas requiring further study.",

			Introduction: "The teacher introduces the revision lesson by asking students to recall the major topics studied during the term.",

			PreviousKnowledge: "Students have studied the major topics covered during the term, including micro-organisms, STIs, plant and animal diseases, food production, preservation and storage.",

			LessonContent: "Revision of micro-organisms; culturing and identification; STIs; activities and control of micro-organisms; classification of plants; plant and animal pests and diseases; disease control; food production; food preservation and storage.",

			TeacherActivities: "The teacher reviews major concepts, asks oral questions, provides revision exercises, explains difficult areas, and guides students through sample examination questions.",

			StudentActivities: "Students answer revision questions, participate in group discussions, solve exercises, ask questions, and identify areas they need to study further.",

			TeachingAids: "Textbooks, past examination questions, revision worksheets, charts, diagrams, whiteboard and marker.",

			Evaluation: "Students complete a comprehensive revision exercise covering all topics taught during the term.",

			Conclusion: "The teacher summarizes the major topics, corrects students' mistakes, provides guidance for further study, and prepares students for the term examination.",

			CreatedBy: adminID,

		},
	}

		for _, lesson_plan := range lessonPlans {
		if err := db.Create(&lesson_plan).Error; err != nil {
			return fmt.Errorf(
				"failed to seed lesson plan %s: %w",
				lesson_plan.BehaviouralObjectives,
				err,
			)
		}

		log.Printf(
			"✅ Seeded lesson plan: Week %d ",
			lesson_plan.BehaviouralObjectives,
		)
	}

	return nil


}