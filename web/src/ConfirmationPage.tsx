import { useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { API_URL } from "./App"

type Status = "idle" | "loading" | "success" | "error"

export const ConfirmationPage = () => {
  const { token = "" } = useParams()
  const redirect = useNavigate()
  const [status, setStatus] = useState<Status>("idle")
  const [errorMsg, setErrorMsg] = useState("")

  const handleConfirm = async () => {
    setStatus("loading")
    try {
      const response = await fetch(`${API_URL}/users/activate/${token}`, {
        method: "PUT",
      })

      if (response.ok) {
        setStatus("success")
        setTimeout(() => redirect("/login"), 2000)
      } else {
        const data = await response.json().catch(() => ({}))
        setErrorMsg(data?.error ?? "Token inválido o expirado.")
        setStatus("error")
      }
    } catch {
      setErrorMsg("Error de conexión. Intenta de nuevo.")
      setStatus("error")
    }
  }

  return (
    <div style={{ maxWidth: 400, margin: "80px auto", textAlign: "center", fontFamily: "sans-serif" }}>
      <h1>Confirmar cuenta</h1>

      {status === "idle" && (
        <>
          <p>Haz clic en el botón para activar tu cuenta.</p>
          <button onClick={handleConfirm} style={btnStyle}>
            Activar cuenta
          </button>
        </>
      )}

      {status === "loading" && <p>Activando tu cuenta...</p>}

      {status === "success" && (
        <p style={{ color: "green" }}>
          ¡Cuenta activada! Redirigiendo al login...
        </p>
      )}

      {status === "error" && (
        <>
          <p style={{ color: "red" }}>{errorMsg}</p>
          <button onClick={handleConfirm} style={btnStyle}>
            Reintentar
          </button>
        </>
      )}
    </div>
  )
}

const btnStyle: React.CSSProperties = {
  padding: "10px 24px",
  fontSize: 16,
  cursor: "pointer",
  borderRadius: 6,
  border: "none",
  background: "#4f46e5",
  color: "#fff",
}
