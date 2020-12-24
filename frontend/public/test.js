"use strict";
const container = document.getElementById("app");
const pokemons = 100;
const fetchData = () => {
    for (let i = 1; i <= pokemons; i++) {
        getPokemon(i);
    }
};
const getPokemon = async (id) => {
    const data = await fetch(`https://pokeapi.co/api/v2/pokemon/${id}`);
    const pokemon = await data.json();
    console.log(pokemon);
    const pokemonType = pokemon.types
        .map((poke) => poke.type.name)
        .join(", ");
    const transformedPokemon = {
        id: pokemon.id,
        name: pokemon.name,
        image: `${pokemon.sprites.front_default}`,
        type: pokemonType
    };
    showPokemon(transformedPokemon);
};
const showPokemon = (pokemon) => {
    let output = `
    <div class="card">
        <span class="card--id">#${pokemon.id}</span>
        <img class="card--image" src=${pokemon.image} alt=${pokemon.name} />
        <h1 class="card--name">${pokemon.name}</h1>
        <span class="card--details">${pokemon.type}</span>
    </div>
    `;
    container.innerHTML += output;
};
fetchData();
