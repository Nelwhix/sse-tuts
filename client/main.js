const evtSource = new EventSource('http://localhost:3000')

evtSource.onmessage = (event) => {
    const data = JSON.parse(event.data)
    const list = document.getElementById('list');
    const newElement = document.createElement("li")
    newElement.textContent = data.data
    list.appendChild(newElement)
}