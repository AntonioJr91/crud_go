const BASE_URL = "http://localhost:8080";

export async function getItems() {
  const resp = await fetch(`${BASE_URL}/items`);
  if (!resp.ok) throw new Error("Erro ao buscar items");
  if(resp.status === 204)return [];
  return  await resp.json();
};

export async function createItem(name, email) {
  const resp = await fetch(`${BASE_URL}/items`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, email }),
  });
  if (!resp.ok) throw new Error("Erro ao criar items");
  return resp.json();
};

export async function getItemById(id) {
  const resp = await fetch(`${BASE_URL}/items/${id}`);
  if (!resp.ok) throw new Error("Item não encontrado");
  return resp.json();
};

export async function updateItem(id, updateData) {
  if (!id) throw new Error("ID inválido!");

  const resp = await fetch(`${BASE_URL}/items/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(updateData)
  });
  if(!resp.ok) throw new Error("Erro ao atualizar item");
  return resp.json();
};

export async function deleteItem(id) {
  const resp = await fetch(`${BASE_URL}/items/${id}`, {
    method: "DELETE",
  });
  if (!resp.ok) throw new Error("Erro ao deletar item");
  if(resp.status === 204) return true;
  alert("Item deletado");
  return resp.json();
};