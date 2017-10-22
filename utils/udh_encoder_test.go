package utils_test

import (
	"testing"
	"utils"

	"github.com/stretchr/testify/assert"
)

func TestInitEncoder(t *testing.T) {
	encoder := utils.InitEncoder()

	assert.NotEmpty(t, encoder)
}

func TestUdhenc_Encode(t *testing.T) {
	encoder := utils.InitEncoder()

	t.Run("GSM 7-bit encode", func(t *testing.T) {
		t.Run("encode regular symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{"48656c6c6f2e2049276d2066696e652c20616e6420796f753f"}}
			m := encoder.Encode("Hello. I'm fine, and you?")
			assert.Equal(t, e, *m)
		})

		t.Run("encode 2 space char symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{"1b2f1b401b281b291b3c1b3e1b3d1b651b14"}}
			m := encoder.Encode(`\|{}[]~€^`)
			assert.Equal(t, e, *m)
		})

		t.Run("encode mixed symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{"4048656c6c6f21205d697d6f2e205e6265722e201315201b2873796d626f6c731b29"}}
			m := encoder.Encode("¡Hello! Ñiño. Über. ΓΩ {symbols}")
			assert.Equal(t, e, *m)
		})

		t.Run("example from MessageBird documentation", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{"546865206d65737361676520746f2062652073656e74"}}
			m := encoder.Encode("The message to be sent")
			assert.Equal(t, e, *m)
		})

		t.Run("mixed symbols splitting", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{
				"546865206d65737361676520746f2062652073656e742c2074686174206e656564732073706c697474696e6720616e642068617320737472616e67652073796d626f6c73206c696b652074686f73653a201b3c1b3e1b281b291b2f2e20486f7765766572206974206e6565647320746f206265206c6f6e676572207468616e206f746865727320736f2069742077696c6c2074616b65206d6f",
				"72652074686174203136302073796d626f6c732e2e2e",
			}}

			m := encoder.Encode(`The message to be sent, that needs splitting and has strange symbols like those: []{}\. However it needs to be longer than others so it will take more that 160 symbols...`)
			assert.Equal(t, e, *m)
		})

		t.Run("splitting for message that is bigger than max SMS capacity", func(t *testing.T) {
			m := encoder.Encode(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. In vehicula aliquam aliquam.
Aliquam mollis nunc quis mi aliquam rutrum. Praesent hendrerit dolor ac ligula suscipit condimentum.
Aenean sollicitudin neque ut ante faucibus porta. Donec eget tortor iaculis, consequat velit lobortis,
iaculis odio. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae;
Integer sed lacus vitae urna malesuada varius. Nunc bibendum elit ex, a elementum eros tempor nec.
Nam quis eros ullamcorper, tincidunt nibh at, aliquet massa. Vivamus pretium maximus tempor.
Vivamus viverra convallis ex. Vivamus est felis, luctus ac urna ut, tristique rhoncus mauris.
Aliquam tincidunt dignissim mi, ut congue neque tincidunt sit amet. Nullam id lorem mauris.
Phasellus luctus arcu at aliquam ornare. Maecenas ornare lorem semper orci auctor, eget placerat sem
finibus. Suspendisse eget auctor sem. Phasellus cursus, orci non blandit tempor, justo neque venenatis
nunc, non facilisis arcu diam ut elit. In at porta mi. Nam nisi risus, commodo eu ex vel, varius fermentum
leo. Curabitur hendrerit nisl massa, eget placerat justo porttitor eget. Integer tempor rhoncus pharetra.
Interdum et malesuada fames ac ante ipsum primis in faucibus. Curabitur vel posuere nunc. Quisque pretium
commodo gravida. Maecenas ultrices metus eu metus commodo, at euismod magna faucibus. Duis volutpat, magna
eu cursus mollis, mauris purus lobortis nulla, vitae sagittis mauris mauris eget tortor. Cras vel ante libero.
Aliquam rhoncus malesuada lacus, sed gravida nisl consequat non. Nunc ut suscipit mi, nec luctus lorem.
Mauris consequat laoreet leo in dignissim. Phasellus et faucibus diam. Proin et ante vulputate elit
finibus laoreet vulputate sed ante. Vivamus blandit eros sed nisl pretium egestas. Ut ac posuere libero,
a rutrum ligula. Pellentesque ac congue nibh.
Etiam elementum aliquet accumsan. Donec auctor porta velit in consectetur. Pellentesque rutrum lacinia
orci ac tempus. In mattis posuere.`)
			assert.Equal(t, 9, len(m.Messages))
		})
	})

	t.Run("UC-2 encode", func(t *testing.T) {
		t.Run("encode UC-2 symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Unicode, []string{"041f044004380432043504420021d83dde00"}}
			m := encoder.Encode("Привет!😀")
			assert.Equal(t, e, *m)
		})

		t.Run("encode only UC-2 symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Unicode, []string{"041f04400438043204350442d83dde000020041a0430043a002004340435043b0430"}}
			m := encoder.Encode("Привет😀 Как дела")
			assert.Equal(t, e, *m)
		})

		t.Run("splitted UC-2 message", func(t *testing.T) {
			e := utils.Encoded{utils.Unicode, []string{
				"042d0442043e00200441043e043e043104490435043d0438043500200434043e043b0436043d043e00200438043c04350442044c00200431043e043b04350435002000370030002004410438043c0432043e043b043e0432002c00200434043b044f00200442043e0433043e002004470442043e0431044b00200431044b",
				"043b043000200432043e0437043c043e0436043d043e04410442044c0020043f0440043e044204350441044204380440043e043204300442044c002004400430043704340435043b0435043d0438043500200441043e043e043104490435043d04380439002e",
			}}

			m := encoder.Encode("Это сообщение должно иметь более 70 символов, для того чтобы была возможность " +
				"протестировать разделение сообщений.")
			assert.Equal(t, e, *m)
		})

		t.Run("splitting for message that is bigger than max SMS capacity", func(t *testing.T) {
			m := encoder.Encode(`Не следует, однако забывать, что начало повседневной работы по формированию
				позиции способствует подготовки и реализации систем массового участия. Не следует, однако забывать,
				что реализация намеченных плановых заданий обеспечивает широкому кругу (специалистов) участие в
				формировании дальнейших направлений развития. Разнообразный и богатый опыт постоянный количественный
				рост и сфера нашей активности способствует подготовки и реализации соответствующий условий активизации.
				Таким образом начало повседневной работы по формированию позиции требуют от нас анализа направлений
				прогрессивного развития. Таким образом сложившаяся структура организации позволяет выполнять важные
				задания по разработке дальнейших направлений развития. Идейные соображения высшего порядка, а также
				постоянный количественный рост и сфера нашей активности играет важную роль в формировании позиций,
				занимаемых участниками в отношении поставленных задач.
				Товарищи! новая модель организационной деятельности играет важную роль в формировании систем массового
				участия. Товарищи! постоянное информационно-пропагандистское обеспечение нашей деятельности представляет
				собой интересный эксперимент проверки позиций, занимаемых участниками в отношении поставленных задач.
				Товарищи! укрепление и развитие структуры способствует подготовки и реализации форм развития.
				Задача организации, в особенности же начало повседневной работы по формированию позиции представляет
				собой интересный эксперимент проверки системы обучения кадров, соответствует насущным потребностям.
				Повседневная практика показывает, что постоянный количественный рост и сфера нашей активности требуют
				от нас анализа дальнейших направлений развития. Равным образом консультация с широким активом
				обеспечивает широкому кругу (специалистов) участие в формировании соответствующий условий активизации.
				Идейные соображения высшего порядка, а также укрепление и развитие структуры требуют определения и
				уточнения системы обучения кадров, соответствует насущным потребностям. Равным образом консультация с
				широким активом влечет за собой процесс внедрения и модернизации позиций, занимаемых участниками в
				отношении поставленных задач. Разнообразный и богатый опыт укрепление и развитие структуры позволяет
				оценить значение существенных финансовых и административных условий.
				С другой стороны дальнейшее развитие различных форм деятельности представляет собой интересный
				эксперимент проверки направлений прогрессивного развития. Таким образом консультация с широким активом
				обеспечивает широкому кругу (специалистов) участие в формировании существенных финансовых и
				административных условий. Равным образом укрепление и развитие структуры представляет собой интересный
				эксперимент проверки новых предложений. Равным образом укрепление и развитие структуры обеспечивает
				широкому кругу (специалистов) участие в формировании систем массового участия.`)
			assert.Equal(t, 9, len(m.Messages))
		})
	})
}

func TestUdhenc_SplitTextMessage(t *testing.T) {
	encoder := utils.InitEncoder()

	t.Run("GSM 7-bit encode", func(t *testing.T) {
		t.Run("encode regular symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{"Hello. I'm fine, and you?"}}
			m := encoder.SplitTextMessage("Hello. I'm fine, and you?")
			assert.Equal(t, e, *m)
		})

		t.Run("encode 2 space char symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{`\|{}[]~€^`}}
			m := encoder.SplitTextMessage(`\|{}[]~€^`)
			assert.Equal(t, e, *m)
		})

		t.Run("encode mixed symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{"¡Hello! Ñiño. Über. ΓΩ {symbols}"}}
			m := encoder.SplitTextMessage("¡Hello! Ñiño. Über. ΓΩ {symbols}")
			assert.Equal(t, e, *m)
		})

		t.Run("example from MessageBird documentation", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{"The message to be sent"}}
			m := encoder.SplitTextMessage("The message to be sent")
			assert.Equal(t, e, *m)
		})

		t.Run("mixed symbols splitting", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{
				`The message to be sent, that needs splitting and has strange symbols like those: []{}\. However it needs to be longer than others so it will take mo`,
				"re that 160 symbols...",
			}}

			m := encoder.SplitTextMessage(`The message to be sent, that needs splitting and has strange symbols like those: []{}\. However it needs to be longer than others so it will take more that 160 symbols...`)
			assert.Equal(t, e, *m)
		})

		t.Run("splitting for message that is bigger than max SMS capacity", func(t *testing.T) {
			m := encoder.SplitTextMessage(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. In vehicula aliquam aliquam.
Aliquam mollis nunc quis mi aliquam rutrum. Praesent hendrerit dolor ac ligula suscipit condimentum.
Aenean sollicitudin neque ut ante faucibus porta. Donec eget tortor iaculis, consequat velit lobortis,
iaculis odio. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae;
Integer sed lacus vitae urna malesuada varius. Nunc bibendum elit ex, a elementum eros tempor nec.
Nam quis eros ullamcorper, tincidunt nibh at, aliquet massa. Vivamus pretium maximus tempor.
Vivamus viverra convallis ex. Vivamus est felis, luctus ac urna ut, tristique rhoncus mauris.
Aliquam tincidunt dignissim mi, ut congue neque tincidunt sit amet. Nullam id lorem mauris.
Phasellus luctus arcu at aliquam ornare. Maecenas ornare lorem semper orci auctor, eget placerat sem
finibus. Suspendisse eget auctor sem. Phasellus cursus, orci non blandit tempor, justo neque venenatis
nunc, non facilisis arcu diam ut elit. In at porta mi. Nam nisi risus, commodo eu ex vel, varius fermentum
leo. Curabitur hendrerit nisl massa, eget placerat justo porttitor eget. Integer tempor rhoncus pharetra.
Interdum et malesuada fames ac ante ipsum primis in faucibus. Curabitur vel posuere nunc. Quisque pretium
commodo gravida. Maecenas ultrices metus eu metus commodo, at euismod magna faucibus. Duis volutpat, magna
eu cursus mollis, mauris purus lobortis nulla, vitae sagittis mauris mauris eget tortor. Cras vel ante libero.
Aliquam rhoncus malesuada lacus, sed gravida nisl consequat non. Nunc ut suscipit mi, nec luctus lorem.
Mauris consequat laoreet leo in dignissim. Phasellus et faucibus diam. Proin et ante vulputate elit
finibus laoreet vulputate sed ante. Vivamus blandit eros sed nisl pretium egestas. Ut ac posuere libero,
a rutrum ligula. Pellentesque ac congue nibh.
Etiam elementum aliquet accumsan. Donec auctor porta velit in consectetur. Pellentesque rutrum lacinia
orci ac tempus. In mattis posuere.`)
			assert.Equal(t, 9, len(m.Messages))
		})
	})

	t.Run("UC-2 encode", func(t *testing.T) {
		t.Run("encode UC-2 symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Unicode, []string{"Привет!😀"}}
			m := encoder.SplitTextMessage("Привет!😀")
			assert.Equal(t, e, *m)
		})

		t.Run("encode only UC-2 symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Unicode, []string{"Привет😀 Как дела"}}
			m := encoder.SplitTextMessage("Привет😀 Как дела")
			assert.Equal(t, e, *m)
		})

		t.Run("splitted UC-2 message", func(t *testing.T) {
			e := utils.Encoded{utils.Unicode, []string{
				"Это сообщение должно иметь более 70 символов, для того чтобы бы",
				"ла возможность протестировать разделение сообщений.",
			}}

			m := encoder.SplitTextMessage("Это сообщение должно иметь более 70 символов, для того чтобы была возможность " +
				"протестировать разделение сообщений.")
			assert.Equal(t, e, *m)
		})

		t.Run("splitting for message that is bigger than max SMS capacity", func(t *testing.T) {
			m := encoder.SplitTextMessage(`Не следует, однако забывать, что начало повседневной работы по формированию
				позиции способствует подготовки и реализации систем массового участия. Не следует, однако забывать,
				что реализация намеченных плановых заданий обеспечивает широкому кругу (специалистов) участие в
				формировании дальнейших направлений развития. Разнообразный и богатый опыт постоянный количественный
				рост и сфера нашей активности способствует подготовки и реализации соответствующий условий активизации.
				Таким образом начало повседневной работы по формированию позиции требуют от нас анализа направлений
				прогрессивного развития. Таким образом сложившаяся структура организации позволяет выполнять важные
				задания по разработке дальнейших направлений развития. Идейные соображения высшего порядка, а также
				постоянный количественный рост и сфера нашей активности играет важную роль в формировании позиций,
				занимаемых участниками в отношении поставленных задач.
				Товарищи! новая модель организационной деятельности играет важную роль в формировании систем массового
				участия. Товарищи! постоянное информационно-пропагандистское обеспечение нашей деятельности представляет
				собой интересный эксперимент проверки позиций, занимаемых участниками в отношении поставленных задач.
				Товарищи! укрепление и развитие структуры способствует подготовки и реализации форм развития.
				Задача организации, в особенности же начало повседневной работы по формированию позиции представляет
				собой интересный эксперимент проверки системы обучения кадров, соответствует насущным потребностям.
				Повседневная практика показывает, что постоянный количественный рост и сфера нашей активности требуют
				от нас анализа дальнейших направлений развития. Равным образом консультация с широким активом
				обеспечивает широкому кругу (специалистов) участие в формировании соответствующий условий активизации.
				Идейные соображения высшего порядка, а также укрепление и развитие структуры требуют определения и
				уточнения системы обучения кадров, соответствует насущным потребностям. Равным образом консультация с
				широким активом влечет за собой процесс внедрения и модернизации позиций, занимаемых участниками в
				отношении поставленных задач. Разнообразный и богатый опыт укрепление и развитие структуры позволяет
				оценить значение существенных финансовых и административных условий.
				С другой стороны дальнейшее развитие различных форм деятельности представляет собой интересный
				эксперимент проверки направлений прогрессивного развития. Таким образом консультация с широким активом
				обеспечивает широкому кругу (специалистов) участие в формировании существенных финансовых и
				административных условий. Равным образом укрепление и развитие структуры представляет собой интересный
				эксперимент проверки новых предложений. Равным образом укрепление и развитие структуры обеспечивает
				широкому кругу (специалистов) участие в формировании систем массового участия.`)
			assert.Equal(t, 9, len(m.Messages))
		})
	})
}

func TestUdhenc_GenerateUDH(t *testing.T) {
	e := utils.InitEncoder()

	t.Run("generates unique UDH for message parts", func(t *testing.T) {
		p := uint8(1)
		parts := uint8(3)
		mesHash := uint32(1)
		udh1 := e.GenerateUDH(p, parts, mesHash)
		p++
		udh2 := e.GenerateUDH(p, parts, mesHash)
		p++
		udh3 := e.GenerateUDH(p, parts, mesHash)

		assert.Equal(t, "050003010301", udh1)
		assert.Equal(t, "050003010302", udh2)
		assert.Equal(t, "050003010303", udh3)
	})

	t.Run("same hash should have same UDH next time", func(t *testing.T) {
		p := uint8(1)
		parts := uint8(3)
		mesHash := uint32(1)
		udh1 := e.GenerateUDH(p, parts, mesHash)
		p++
		udh2 := e.GenerateUDH(p, parts, mesHash)
		p++
		udh3 := e.GenerateUDH(p, parts, mesHash)

		assert.Equal(t, "050003010301", udh1)
		assert.Equal(t, "050003010302", udh2)
		assert.Equal(t, "050003010303", udh3)
	})
}
