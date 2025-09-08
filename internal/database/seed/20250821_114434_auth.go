package seed

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func Auth() {
	app.PrintInfo("seeding auth...")
	// acciones o verbos de autorización
	var actions = map[string][]string{
		"crud": {
			"view",   // ver un recurso
			"create", // crear un recurso
			"update", // modificar un recurso
			"delete", // eliminar un recurso
		},

		"moderator": {
			"approve", // aprobar contenido pendiente
			"reject",  // rechazar contenido reportado o pendiente
			"publish", // publicar contenido en nombre de la comunidad
			"lock",    // cerrar un hilo, comentario o bloquear contenido
			"warn",    // advertir a un usuario
			"mute",    // silenciar temporalmente a un usuario
			"suspend", // suspender la cuenta de un usuario
			"ban",     // bloquear permanentemente a un usuario
			"kick",    // expulsar a un usuario de un grupo/canal
			"enable",  // reactivar un recurso
			"disable", // desactivar un recurso
			"archive", // archivar contenido viejo o inactivo
			"restore", // restaurar contenido archivado
			"import",  // importar datos o recursos
			"export",  // exportar datos o recursos
		},

		"user": {
			"comment",   // escribir un comentario
			"reply",     // responder a un comentario
			"mention",   // mencionar a otro usuario (@user)
			"tag",       // etiquetar un recurso o persona
			"share",     // compartir un recurso
			"like",      // dar un me gusta
			"react",     // reaccionar con un emoji u otra acción
			"follow",    // seguir a un usuario o recurso
			"invite",    // invitar a alguien a la comunidad
			"join",      // unirse a un grupo o canal
			"report",    // reportar contenido o usuario
			"vote",      // votar en encuestas o contenido
			"subscribe", // suscribirse a un recurso o canal
			"bookmark",  // guardar como favorito
			"highlight", // resaltar un texto o recurso
			"pin",       // fijar contenido importante
		},

		"admin": {
			"grant permission",  // otorgar permisos a otros usuarios
			"revoke permission", // revocar permisos a otros usuarios
			"grant role",        // otorgar roles a otros usuarios
			"revoke role",       // revocar roles a otros usuarios
		},
	}
	// slice con los nombres de los modelos
	models := []string{"user", "role", "permission", "city", "state", "country"}

	// slice con las referencias de los permisos
	adminPermissions := []bson.ObjectID{}
	crudPermissions := []bson.ObjectID{}
	moderatorPermissions := []bson.ObjectID{}
	userPermissions := []bson.ObjectID{}

	// rol de usuario
	for _, action := range actions["user"] {
		permission := model.NewPermission()
		permission.Name = action
		if err := permission.Create(); err != nil {
			app.PrintError("Fail to create permission: :permission :error", app.E("permission", permission.Name), app.E("error", err.Error()))
			panic(err)
		}
		userPermissions = append(userPermissions, permission.ID)
	}
	roleUser := model.NewRole()
	roleUser.Name = "user"
	roleUser.PermissionIDs = userPermissions
	if err := roleUser.Create(); err != nil {
		app.PrintError("Fail to create role: :role :error", app.E("role", roleUser.Name), app.E("error", err.Error()))
		panic(err)
	}

	// rol de moderador
	for _, action := range actions["moderator"] {
		permission := model.NewPermission()
		permission.Name = action
		if err := permission.Create(); err != nil {
			app.PrintError("Fail to create permission: :permission :error", app.E("permission", permission.Name), app.E("error", err.Error()))
			panic(err)
		}
		moderatorPermissions = append(moderatorPermissions, permission.ID)
	}
	roleModerator := model.NewRole()
	roleModerator.Name = "moderator"
	roleModerator.PermissionIDs = moderatorPermissions
	if err := roleModerator.Create(); err != nil {
		app.PrintError("Fail to create role: :role :error", app.E("role", roleModerator.Name), app.E("error", err.Error()))
		panic(err)
	}

	// rol de admin
	for _, m := range models {
		for _, action := range actions["crud"] {
			permission := model.NewPermission()
			permission.Name = action + " " + m
			if err := permission.Create(); err != nil {
				app.PrintError("Fail to create permission: :permission :error", app.E("permission", permission.Name), app.E("error", err.Error()))
				panic(err)
			}
			crudPermissions = append(crudPermissions, permission.ID)
		}
	}

	for _, action := range actions["admin"] {
		permission := model.NewPermission()
		permission.Name = action
		if err := permission.Create(); err != nil {
			app.PrintError("Fail to create permission: :permission :error", app.E("permission", permission.Name), app.E("error", err.Error()))
			panic(err)
		}
		adminPermissions = append(adminPermissions, permission.ID)
	}

	adminPermissions = append(adminPermissions, crudPermissions...)
	adminPermissions = append(adminPermissions, moderatorPermissions...)
	adminPermissions = append(adminPermissions, userPermissions...)

	roleAdmin := model.NewRole()
	roleAdmin.Name = "admin"
	roleAdmin.PermissionIDs = adminPermissions
	if err := roleAdmin.Create(); err != nil {
		app.PrintError("Fail to create role: :role :error", app.E("role", roleAdmin.Name), app.E("error", err.Error()))
		panic(err)
	}

	// usuario admin
	userAdmin := model.NewUser()
	city := model.NewCity()
	if err := city.First("name", "Medellín"); err != nil {
		app.PrintError("Fail to find city: :city :error", app.E("city", city.Name), app.E("error", err.Error()))
		panic(err)
	}
	hashedPassword, er := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if er != nil {
		app.PrintError("Fail to generate password: :error", app.E("error", er.Error()))
		panic(er)
	}
	userAdmin.Email = "admin@gmail.com"
	userAdmin.Password = string(hashedPassword)
	userAdmin.Profile = &model.Profile{
		Nickname:        "admin",
		FullName:        "admin",
		PhoneNumber:     "+573203110099",
		DiscordUsername: "admin",
		CityID:          city.ID,
	}
	userAdmin.RoleIDs = []bson.ObjectID{roleAdmin.ID}
	if err := userAdmin.Create(); err != nil {
		app.PrintError("Fail to create user: :user :error", app.E("user", userAdmin.Email), app.E("error", err.Error()))
		panic(err)
	}

	app.PrintInfo("Finish seed auth")

}
