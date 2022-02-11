<script setup>
  import { useGet } from "./hooks/useHttp";

  const notes = useGet("/api/note");
  console.log(notes)

</script>

<template>
  <div id="app">
    <li v-for="note in notes">
        {{note.author}} {{note.content}} {{note.created_at}}
    </li>
  </div>
</template>
<style>
#app {
  font-family: Helvetica;
  font-size: 36px;
}
</style>
