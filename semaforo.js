// semaforo.js
// Lógica para cambiar los colores del semáforo en un bucle continuo
function cambiarColor() {
  const semaforo = document.querySelector(".verde");
  setInterval(() => {
    semaforo.classList.toggle("verde");
    semaforo.classList.toggle("rojo");
  }, 800); // Cambio de colores cada 10 segundos
}

// Iniciamos la simulación
cambiarColor();
