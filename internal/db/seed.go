package db

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/felipeeguia03/vol7/internal/store"
)

var usernames = []string{
	"byteBandit",
	"golangGuru",
	"asyncAce",
	"codeNomad",
	"devDrifter",
	"stackSamurai",
	"nilNavigator",
	"panicPilot",
	"mutexMaster",
	"sliceSurfer",
	"structSmith",
	"interfaceIvy",
	"goroutineGuy",
	"channelChaser",
	"pointerPro",
	"loopLegend",
	"binaryBlaze",
	"cacheCrafter",
	"heapHero",
	"mapMaven",
	"funcFury",
	"compileKing",
	"debugDynamo",
	"refactorRex",
	"testTitan",
	"deployDuke",
	"cloudCoder",
	"lambdaLion",
	"dataDrake",
	"queryQueen",
	"sqlSage",
	"jsonJedi",
	"apiApex",
	"cryptoCliff",
	"hashHawk",
	"packetPioneer",
	"runtimeRider",
	"threadThor",
	"bitBuilder",
	"zeroZen",
	"syntaxSonic",
	"vectorVince",
	"kernelKnight",
	"byteBison",
	"protoPulse",
	"scriptStorm",
	"nodeNinja",
	"fiberFalcon",
	"logicLuna",
	"atomicArrow",
}

var titles = []string{
	"Construí esto en Go y no esperaba que fuera tan rápido",
	"Hoy entendí algo clave sobre arquitectura limpia",
	"El bug que me enseñó más que 10 tutoriales",
	"De junior a backend sólido: lo que realmente importa",
	"Optimistic locking explicado sin humo",
	"Si tu API es lenta, probablemente sea esto",
	"Lo que nadie te dice sobre aprender a programar",
	"Context en Go: la pieza que cambia todo",
	"El error que casi rompe producción",
	"Microservicios no son la solución mágica",
	"Cómo estructuro mis proyectos en Go",
	"El poder real de los middlewares",
	"Aprender leyendo código > viendo cursos",
	"Por qué tus handlers están demasiado gordos",
	"La diferencia entre saber código y saber diseñar",
	"El día que entendí concurrencia de verdad",
	"Tu base de datos también tiene sentimientos",
	"Versionado: el héroe silencioso de la concurrencia",
	"Menos frameworks, más fundamentos",
	"Debuggear también es ingeniería",
	"Cómo pienso antes de escribir una línea de código",
	"Refactorizar no es opcional",
	"Lo simple escala mejor",
	"El problema no era el lenguaje",
	"Cómo evitar lost updates en producción",
	"El patrón que me ahorró horas de soporte",
	"Escribir menos código es ganar",
	"Tu arquitectura habla de vos",
	"Separación de responsabilidades en la práctica",
	"El backend no es solo CRUD",
	"Lo que aprendí rompiendo cosas",
	"Timeouts: el detalle que evita incendios",
	"Tu API necesita límites",
	"Interfaces bien usadas cambian el diseño",
	"Programar es tomar decisiones",
	"El costo oculto de no testear",
	"Concurrencia sin miedo",
	"Tu código es tu reputación",
	"El arte de decir 'no' a una feature",
	"Cuando todo parece funcionar… pero no escala",
	"Pequeños cambios, gran impacto",
	"El problema no era el bug, era el diseño",
	"Menos magia, más claridad",
	"Por qué amo Go para backend",
	"Arquitectura > sintaxis",
	"Errores que me hicieron mejor ingeniero",
	"El secreto no es velocidad, es consistencia",
	"Middleware bien usado = código limpio",
	"Construir criterio es más difícil que aprender sintaxis",
	"Hoy optimicé algo que nadie veía",
}

var content = []string{
	"Pensé que iba a necesitar optimizar todo. Al final, el 80% del rendimiento vino de elegir bien la estructura de datos. Go hace el resto.",
	"Arquitectura limpia no es sobre carpetas bonitas. Es sobre dependencias que apuntan hacia adentro, no hacia afuera.",
	"Un slice mal indexado me explotó en producción. Dolió. Pero ahora valido todo lo que entra.",
	"No es aprender más frameworks. Es entender HTTP, concurrencia y bases de datos. Lo demás es accesorio.",
	"No bloquees filas si no hace falta. Agregá una columna version y dormí tranquilo.",
	"Antes de escalar horizontalmente, mirá tus queries. Muchas veces el cuello de botella está en un SELECT mal indexado.",
	"Tutoriales enseñan cómo escribir código. Los errores te enseñan cuándo no hacerlo.",
	"Context no es magia. Es cancelación, timeout y propagación de request. Si lo entendés, tu backend mejora.",
	"El error no era grave. El problema fue no tener logging suficiente para entenderlo rápido.",
	"Microservicios sin necesidad real = complejidad distribuida. Primero entendé tu dominio.",
	"main pequeño, dependencias claras, handlers finos, store aislado. Simplicidad gana.",
	"Middleware no es para lógica de negocio. Es para comportamiento transversal.",
	"Leer código open source te enseña decisiones reales, no ejemplos de laboratorio.",
	"Si tu handler tiene 200 líneas, probablemente esté haciendo demasiado.",
	"Saber sintaxis te hace programador. Saber diseño te hace ingeniero.",
	"Concurrencia no es paralelismo. Es coordinación. Y eso cambia todo.",
	"Un index bien puesto puede ser más poderoso que cualquier optimización.",
	"Version = version + 1 puede salvarte de perder datos silenciosamente.",
	"Cada dependencia extra es una deuda futura.",
	"Debuggear no es suerte. Es método.",
	"Antes de escribir código, escribí qué problema estás resolviendo.",
	"Refactorizar es pagar deuda técnica antes de que cobre intereses.",
	"Lo que no entendés no escala.",
	"No era el lenguaje. Era la falta de diseño.",
	"Lost updates no hacen ruido. Por eso son peligrosos.",
	"Un buen patrón elimina código repetido sin esconder la intención.",
	"Más código = más superficie de bugs.",
	"La arquitectura es la conversación silenciosa entre tus capas.",
	"Separar responsabilidades reduce miedo al cambiar código.",
	"CRUD es el inicio. El valor está en las reglas de negocio.",
	"Romper producción una vez te enseña más que cien deploys exitosos.",
	"Timeouts no son opcionales. Son protección.",
	"Sin límites claros, tu API es vulnerable.",
	"Una interfaz pequeña bien definida vale más que cinco grandes mal pensadas.",
	"Programar es elegir restricciones correctas.",
	"Tests no previenen bugs. Previenen miedo a cambiar.",
	"Concurrencia sin control es caos elegante.",
	"El código que escribís hoy es el que otro va a mantener mañana.",
	"No toda feature merece existir.",
	"Escalar es más difícil que funcionar.",
	"A veces el mayor impacto es eliminar código.",
	"Un bug repetido suele ser síntoma de mal diseño.",
	"La claridad es más importante que la inteligencia.",
	"Go no es mágico. Es predecible. Y eso es poderoso.",
	"Arquitectura es decidir dónde vive cada responsabilidad.",
	"Cada error es feedback técnico.",
	"La consistencia construye sistemas confiables.",
	"Middleware bien usado elimina duplicación sin ocultar intención.",
	"Aprender sintaxis es rápido. Desarrollar criterio toma años.",
	"Optimizar algo invisible puede mejorar todo el sistema.",
}

var tags = []string{
	"golang",
	"backend",
	"webdevelopment",
	"programming",
	"softwareengineering",
	"cleanarchitecture",
	"microservices",
	"api",
	"rest",
	"concurrency",
	"databases",
	"postgres",
	"devlife",
	"coding",
	"cloud",
	"scalability",
	"optimisticlocking",
	"systemdesign",
	"techmindset",
	"buildinpublic",
}

var comments = []string{
	"Gran punto. Mucha gente subestima lo importante que es el diseño antes del código.",
	"Totalmente de acuerdo, la arquitectura es lo que realmente escala.",
	"Esto me pasó hace poco en producción, aprendí por las malas",
	"Excelente explicación, especialmente la parte de concurrencia.",
	"Go realmente te obliga a pensar en simplicidad.",
	"Me encantó el enfoque práctico, más contenido así.",
	"El tema de optimistic locking casi nadie lo menciona.",
	"Esto debería enseñarse más en lugar de solo frameworks.",
	"Gran recordatorio sobre mantener handlers pequeños.",
	"La claridad siempre gana sobre la complejidad innecesaria.",
	"Justo estoy trabajando en algo similar, me sirve mucho.",
	"Middleware bien usado cambia completamente la estructura del proyecto.",
	"Qué buena reflexión sobre criterio vs sintaxis.",
	"Esto aplica incluso fuera de Go, es puro diseño.",
	"Excelente hilo, se nota la experiencia detrás.",
	"Tal cual, menos magia y más fundamentos.",
	"Esto explica muchos bugs silenciosos que he visto.",
	"Me hizo replantear cómo estoy organizando mis servicios.",
	"Necesitamos más contenido técnico de este nivel.",
	"Gran aporte, directo y sin humo.",
}

func Seed(store store.Storage, db *sql.DB) error {
	ctx := context.Background()

	users := generateUser(100)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin users transaction: %w", err)
	}

	for _, user := range users {
		err := store.Users.Create(ctx, tx, user)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error creando usuario: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit users transaction: %w", err)
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		err := store.Posts.Create(ctx, post)
		if err != nil {
			fmt.Println("error creando post", err)
		}
	}

	comments := generateComments(100, users, posts)
	for _, cmt := range comments {
		err := store.Comments.Create(ctx, cmt)
		if err != nil {
			fmt.Println("error creando comentario", err)
		}
	}

	return nil
}

func generateUser(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "hotmail.com",
		}
		users[i].Password.Set("123")
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: content[rand.Intn(len(content))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
			Comments: []*store.Comment{},
		}
	}
	return posts

}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cmts := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]

		cmts[i] = &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: content[rand.Intn(len(content))],
		}

	}
	return cmts
}
