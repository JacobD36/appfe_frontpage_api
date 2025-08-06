package template

const (
	// Base template HTML5 con estilos comunes
	baseHTMLTemplate = `<!DOCTYPE html>
<html lang="es">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	<style>
		body { margin: 0; padding: 0; font-family: 'Roboto', Helvetica, Arial, sans-serif; }
		.header { background-color: #185691; color: white; font-size: 2rem; padding: 10px; text-align: center; }
		.content { text-align: center; padding: 20px; }
		.message { font-size: 1.2rem; color: #3c3c3c; text-align: justify; line-height: 1.6; }
		.highlight-box { text-align: center; margin: 20px 0; }
		.highlight-label { font-size: 1.1rem; color: #185691; font-weight: bold; }
		.highlight-value { 
			font-size: 1.3rem; 
			font-weight: bold; 
			color: #d32f2f; 
			background-color: #f5f5f5; 
			padding: 8px 16px; 
			border-radius: 4px; 
			display: inline-block; 
			margin-top: 8px; 
			font-family: 'Courier New', monospace; 
		}
		.button {
			display: inline-block;
			background-color: #185691;
			color: white;
			text-decoration: none;
			padding: 12px 24px;
			border-radius: 5px;
			font-weight: bold;
			margin: 20px 0;
		}
		.footer { 
			background-color: #4a4a4a; 
			color: white; 
			text-align: center; 
			padding: 10px 0; 
			margin-top: 30px; 
		}
	</style>
</head>
<body>
	<div class="header">%s</div>
	<div class="content">
		%s
	</div>
	<div class="footer">
		GeekWay 2025 - Todos los derechos reservados
	</div>
</body>
</html>`

	// Template específico para email de bienvenida
	welcomeContentTemplate = `
		<div class="message">
			¡Hola %s! <br><br>
			
			Te damos la bienvenida al portal administrativo de APPFE Lima. Tu cuenta ha sido creada exitosamente y ya puedes comenzar a usar nuestra plataforma.
			<br><br>
			
			<div class="highlight-box">
				<div class="highlight-label">Tu contraseña de acceso:</div>
				<div class="highlight-value">%s</div>
			</div>
			
			Si tienes alguna pregunta, no dudes en contactarnos.
		</div>`

	// Template para reset de contraseña
	passwordResetContentTemplate = `
		<div class="message">
			Hola %s, <br><br>
			
			Recibimos una solicitud para restablecer tu contraseña. Si fuiste tú quien la solicitó, haz clic en el siguiente enlace:
			<br><br>
			
			<div class="highlight-box">
				<a href="%s" class="button">Restablecer Contraseña</a>
			</div>
			
			Este enlace expirará en 1 hora por seguridad.<br><br>
			
			Si no solicitaste este cambio, puedes ignorar este correo.
		</div>`

	// Template para validación de email
	emailValidationContentTemplate = `
		<div class="message">
			¡Hola %s! <br><br>
			
			Para completar tu registro en APPFE Lima, necesitamos verificar tu dirección de correo electrónico.
			<br><br>
			
			<div class="highlight-box">
				<a href="%s" class="button">Verificar Email</a>
			</div>
			
			Este enlace expirará en 24 horas.<br><br>
			
			Si no te registraste en nuestra plataforma, puedes ignorar este correo.
		</div>`
)
