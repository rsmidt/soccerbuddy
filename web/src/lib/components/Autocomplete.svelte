<script lang="ts">
  import type { AutocompleteOption } from "$lib/components/autocomplete";
  import type { HTMLInputAttributes } from "svelte/elements";

  type InputProps = Omit<HTMLInputAttributes, "type" | "value" | "oninput" | "onkeydown">;

  const {
    fetchOptions,
    select,
    options = [],
    ...rest
  }: {
    fetchOptions: (query: string) => void;
    select: (option: AutocompleteOption) => void;
    options: AutocompleteOption[];
  } & InputProps = $props();

  let inputValue = $state("");
  let showList = $state(false);
  let highlightedIndex = $state(-1);

  // Click outside action to close the suggestions list.
  function clickOutside(node: HTMLElement, callback: () => void) {
    $effect(() => {
      const handleClick = (event: MouseEvent) => {
        if (!node.contains(event.target as Node)) {
          callback();
        }
      };
      document.addEventListener("click", handleClick);
      return () => {
        document.removeEventListener("click", handleClick);
      };
    });
  }

  // Handler for input events.
  function handleInput(event: Event) {
    const target = event.target as HTMLInputElement;
    inputValue = target.value;
    if (inputValue.trim().length > 0) {
      fetchOptions(inputValue);
      showList = true;
      highlightedIndex = -1;
    } else {
      showList = false;
    }
  }

  // Function to handle option selection.
  function selectOption(option: AutocompleteOption) {
    select(option);
    inputValue = "";
    showList = false;
    highlightedIndex = -1;
  }

  // Close the list when clicking outside.
  function handleClickOutside() {
    showList = false;
  }

  // Handle Keydown Events for Navigation and Selection.
  function handleKeyDown(event: KeyboardEvent) {
    if (showList && options.length > 0) {
      if (event.key === "ArrowDown") {
        // Move highlight down.
        highlightedIndex = (highlightedIndex + 1) % options.length;
        event.preventDefault();
      } else if (event.key === "ArrowUp") {
        // Move highlight up.
        highlightedIndex = (highlightedIndex - 1 + options.length) % options.length;
        event.preventDefault();
      } else if (event.key === "Enter") {
        // Select the highlighted option.
        if (highlightedIndex >= 0 && highlightedIndex < options.length) {
          selectOption(options[highlightedIndex]);
          event.preventDefault();
        }
      } else if (event.key === "Escape") {
        // Close the suggestion list.
        showList = false;
      }
    }
  }
</script>

<div class="autocomplete" use:clickOutside={handleClickOutside}>
  <input
    {...rest}
    type="text"
    bind:value={inputValue}
    oninput={handleInput}
    onkeydown={handleKeyDown}
    aria-autocomplete="list"
    aria-expanded={showList}
    aria-controls="autocomplete-list"
    aria-activedescendant={highlightedIndex >= 0 ? `option-${highlightedIndex}` : undefined}
  />
  {#if showList && options.length > 0}
    <ul id="autocomplete-list" class="options" role="listbox">
      {#each options as option, index}
        {@const isSelected = index === highlightedIndex}
        <li
          id={`option-${index}`}
          role="option"
          class:selected={isSelected}
          onclick={() => selectOption(option)}
          onkeydown={() => {}}
          onmouseenter={() => (highlightedIndex = index)}
          aria-selected={isSelected}
        >
          {option.name}
        </li>
      {/each}
    </ul>
  {/if}
</div>

<style>
  .autocomplete {
    position: relative;
    width: 100%;
    max-width: 400px;
    margin: 0 auto; /* Center on larger screens */
  }

  input {
    width: 100%;
  }

  .options {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    border: 1px solid hsl(var(--gray-300));
    border-top: none;
    max-height: 200px;
    overflow-y: auto;
    background-color: var(--bg-100);
    z-index: 1000;
    list-style: none;
    margin: 0;
    padding: 0;
    border-radius: 0 0 0.25rem 0.25rem;
  }

  .options li {
    padding: 0.75rem 1rem;
    cursor: pointer;
    color: var(--text-100);
    font-size: var(--text-md);
  }

  .options li:hover,
  .options li.selected {
    background-color: hsl(var(--primary-500));
    color: var(--text-100);
  }

  /* Mobile-first adjustments */
  @media (min-width: 600px) {
    .autocomplete {
      max-width: 600px;
    }

    input {
      font-size: var(--text-lg);
      padding: 1rem 1.25rem;
    }

    .options li {
      font-size: var(--text-lg);
      padding: 1rem 1.25rem;
    }
  }
</style>
