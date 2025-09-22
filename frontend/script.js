import { getItems, createItem, getItemById, updateItem, deleteItem } from "./api.js";
const app = document.getElementById('app');
var cards = ["Adicionar", "Listar Todos", "Buscar por ID", "Editar", "Deletar"];

function renderHome() {
  app.innerHTML = `
    <h1 class="text-2xl font-bold text-center mb-6">Menu Principal</h1>
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-6 justify-items-center">
      ${cards.map(action => `
        <div onclick="handleCardClick('${action}')" 
             class="bg-white shadow-md rounded-2xl p-6 w-48 text-center cursor-pointer hover:bg-gray-50">
          <p class="font-semibold text-gray-700">${action}</p>
        </div>
      `).join('')}
    </div>
  `;
};

function handleCardClick(action) {
  switch (action) {
    case 'Adicionar': renderAddForm(); break;
    case 'Listar todos': renderListAll(); break;
    case 'Buscar por ID': renderSearchById(); break;
    case 'Editar': renderEditItem(); break;
    case 'Deletar': renderDeleteItem(); break;
  }
};
window.handleCardClick = handleCardClick;

function renderAddForm() {
  app.innerHTML = `
    <h2 class="text-xl font-bold mb-4 text-center">Adicionar Item</h2>
    <form class="bg-white shadow-md rounded-xl p-6" onsubmit="addItem(event)">
      <input id="name" type="text" placeholder="Nome" class="border rounded w-full p-2 mb-4" required />
      <input id="email" type="email" placeholder="Email" class="border rounded w-full p-2 mb-4" required />
      <button class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 w-full">Cadastrar</button>
    </form>
    <button onclick="renderHome()" class="mt-4 w-full bg-gray-300 px-4 py-2 rounded hover:bg-gray-400">Voltar</button>
  `;
};

async function addItem(event) {
  event.preventDefault();
  const name = document.getElementById('name').value;
  const email = document.getElementById('email').value;
  try {
    await createItem(name, email);
    alert('Item cadastrado com sucesso!');
    renderHome();
  } catch (error) {
    console.log(error.message);
  }
};
window.addItem = addItem;

async function renderListAll() {
  try {
    const items = await getItems();
      app.innerHTML = `
        <p class="text-center text-gray-600">Nenhum item cadastrado.</p>
        <button onclick="renderHome()" class="mt-4 w-full bg-gray-300 px-4 py-2 rounded hover:bg-gray-400">Voltar</button>
      `;

    const listHTML = items
      .sort((a, b) => a.id - b.id)
      .map(item => `
        <li class="border-b py-2">ID: ${item.id} - <strong>${item.name}</strong> (${item.email})</li>
      `).join('');

    app.innerHTML = `
      <h2 class="text-xl font-bold mb-4 text-center">Lista de Itens</h2>
      <ul class="bg-white shadow-md rounded-xl p-4 max-h-80 overflow-y-auto">${listHTML}</ul>
      <button onclick="renderHome()" class="mt-4 w-full bg-gray-300 px-4 py-2 rounded hover:bg-gray-400">Voltar</button>
    `;
  } catch (error) {
    console.log(error.message);
  }

};

function renderSearchById() {
  app.innerHTML = `
    <h2 class="text-xl font-bold mb-4 text-center">Buscar por ID</h2>
    <form class="bg-white shadow-md rounded-xl p-6" onsubmit="searchById(event)">
      <input id="searchId" type="number" placeholder="Digite o ID" class="border rounded w-full p-2 mb-4" required />
      <button class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 w-full">Buscar</button>
    </form>
    <button onclick="renderHome()" class="mt-4 w-full bg-gray-300 px-4 py-2 rounded hover:bg-gray-400">Voltar</button>
  `;
};

async function searchById(event) {
  event.preventDefault();
  const id = parseInt(document.getElementById('searchId').value);
  try {
    const item = await getItemById(id);
       app.innerHTML = `
      <h2 class="text-xl font-bold mb-4 text-center">Lista de Itens</h2>
      <ul class="bg-white shadow-md rounded-xl p-4 max-h-80 overflow-y-auto">${`ID: ${item.id} - <strong>${item.name}</strong> (${item.email}`}</ul>
      <button onclick="renderHome()" class="mt-4 w-full bg-gray-300 px-4 py-2 rounded hover:bg-gray-400">Voltar</button>
    `;
  } catch (error) {
    alert("Item não encontrado");
    console.log(error.message);
  }
};
window.searchById = searchById;

function renderEditItem() {
  app.innerHTML = `
    <h2 class="text-xl font-bold mb-4 text-center">Editar Item</h2>
    <form class="bg-white shadow-md rounded-xl p-6" onsubmit="editItemHanlder(event)">
      <input id="editId" type="number" placeholder="ID do item" class="border rounded w-full p-2 mb-4" required />
      <input id="editName" type="text" placeholder="Novo Nome" class="border rounded w-full p-2 mb-4" required />
      <input id="editEmail" type="email" placeholder="Novo Email" class="border rounded w-full p-2 mb-4" required />
      <button class="bg-yellow-500 text-white px-4 py-2 rounded hover:bg-yellow-600 w-full">Salvar</button>
    </form>
    <button onclick="renderHome()" class="mt-4 w-full bg-gray-300 px-4 py-2 rounded hover:bg-gray-400">Voltar</button>
  `;
};

async function editItemHanlder(event) {
  event.preventDefault();
  const id = parseInt(document.getElementById('editId').value);
  const name = document.getElementById('editName').value;
  const email = document.getElementById('editEmail').value;

  try {
    await updateItem(id, {name, email});
    alert("Item atualizado com sucesso!");
    renderHome();
  } catch (error) {
    alert("Erro ao atualizar.");
    console.log(error.message);
  }
};
window.editItemHanlder = editItemHanlder;

function renderDeleteItem() {
  app.innerHTML = `
    <h2 class="text-xl font-bold mb-4 text-center">Deletar Item</h2>
    <form class="bg-white shadow-md rounded-xl p-6" onsubmit="deleteItemHandler(event)">
      <input id="deleteId" type="number" placeholder="ID do item" class="border rounded w-full p-2 mb-4" required />
      <button class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 w-full">Deletar</button>
    </form>
    <button onclick="renderHome()" class="mt-4 w-full bg-gray-300 px-4 py-2 rounded hover:bg-gray-400">Voltar</button>
  `;
};

async function deleteItemHandler(event) {
  event.preventDefault();
  const id = parseInt(document.getElementById('deleteId').value);

  try {
    await deleteItem(id);
    alert("Item deletado com sucesso!");
    renderHome();
  } catch (error) {
    alert("Erro ao deletar");
    console.log(error.message);
  }
};
window.deleteItemHandler = deleteItemHandler;

renderHome();
window.renderHome = renderHome;
