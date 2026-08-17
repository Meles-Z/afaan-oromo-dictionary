<script lang="ts">
  import { SearchWords } from '../wailsjs/go/main/App';
  import type { dictionary } from '../wailsjs/go/models';
  import { onMount } from 'svelte';


  let inputEl: HTMLInputElement;

  onMount(() => {
    inputEl?.focus();
  });

  type Direction = 'en' | 'om';

  let query = '';
  let results: dictionary.Word[] = [];
  let loading = false;
  let searched = false;
  let direction: Direction = 'en';
  let swapping = false;
  let selectedWord: dictionary.Word | null = null;

  let debounceTimer: ReturnType<typeof setTimeout>;
  let requestId = 0;

  const labels = {
    en: { name: 'English', placeholder: 'Type an English word…' },
    om: { name: 'Afaan Oromo', placeholder: 'Barreessi jecha Afaan Oromoo…' }
  };

  // Capitalizes only the first letter of the whole string, leaves the rest untouched
  function cap(text: string): string {
    if (!text) return text;
    return text.charAt(0).toUpperCase() + text.slice(1);
  }

  function onInput() {
    clearTimeout(debounceTimer);
    const trimmed = query.trim();
    selectedWord = null;

    if (trimmed === '') {
      results = [];
      loading = false;
      searched = false;
      return;
    }

    debounceTimer = setTimeout(() => handleSearch(trimmed), 200);
  }

  async function handleSearch(q: string) {
    const thisRequest = ++requestId;
    loading = true;

    try {
      const res = await SearchWords(q, direction);
      if (thisRequest !== requestId) return;
      results = res ?? [];
      searched = true;
    } catch (err) {
      if (thisRequest !== requestId) return;
      console.error('Search failed:', err);
      results = [];
      searched = true;
    } finally {
      if (thisRequest === requestId) loading = false;
    }
  }

  function swapDirection() {
    swapping = true;
    direction = direction === 'en' ? 'om' : 'en';
    query = '';
    results = [];
    searched = false;
    loading = false;
    selectedWord = null;
    setTimeout(() => (swapping = false), 400);
  }

  function selectDirection(d: Direction) {
    if (d === direction) return;
    swapDirection();
  }

  function openDetail(word: dictionary.Word) {
    selectedWord = word;
  }

  function backToResults() {
    selectedWord = null;
  }
</script>

<main>
  <header>
    <h1>Kitaaba Jechootaa</h1>
    <p class="subtitle">English ↔ Afaan Oromo dictionary</p>
  </header>

  <div class="translator">
    <!-- Language selector row -->
    <div class="lang-row">
      <button class="lang-pill" class:active={direction === 'en'} on:click={() => selectDirection('en')}>
        English
      </button>

      <button class="swap-btn" class:spinning={swapping} on:click={swapDirection} aria-label="Swap languages" title="Swap languages">
        <svg viewBox="0 0 24 24" width="20" height="20">
          <path fill="currentColor" d="M7 7h11l-3.5-3.5L16 2l6 6-6 6-1.5-1.5L18 9H7V7zm10 10H6l3.5 3.5L8 22l-6-6 6-6 1.5 1.5L6 15h11v2z"/>
        </svg>
      </button>

      <button class="lang-pill" class:active={direction === 'om'} on:click={() => selectDirection('om')}>
        Afaan Oromo
      </button>
    </div>

    <!-- Search input -->
    <div class="search-box">
      <input
        bind:this={inputEl}
        type="text"
        bind:value={query}
        on:input={onInput}
        placeholder={labels[direction].placeholder}
        class:oromo={direction === 'om'}
      />
      {#if loading}
        <span class="spinner" aria-hidden="true"></span>
      {/if}
    </div>

    {#if selectedWord}
      <!-- LARGE DETAIL VIEW: replaces the list entirely while active -->
      <div class="detail-large">
        <button class="back-btn" on:click={backToResults}>
          ← Back to results
        </button>

        <div class="detail-term-large">
          {cap(direction === 'en' ? selectedWord.english : selectedWord.afaanOromo)}
          {#if selectedWord.partOfSpeech}
            <span class="detail-pos-large">{selectedWord.partOfSpeech}</span>
          {/if}
        </div>

        <div class="detail-translation-large">
          {cap(direction === 'en' ? selectedWord.afaanOromo : selectedWord.english)}
        </div>

        {#if selectedWord.pronunciation}
          <div class="detail-pron-large">/{selectedWord.pronunciation}/</div>
        {/if}

        {#if selectedWord.exampleEn || selectedWord.exampleOm}
          <div class="detail-examples-large">
            {#if selectedWord.exampleEn}
              <p class="ex-en">{selectedWord.exampleEn}</p>
            {/if}
            {#if selectedWord.exampleOm}
              <p class="ex-om">{selectedWord.exampleOm}</p>
            {/if}
          </div>
        {/if}
      </div>
    {:else}
      <!-- RESULTS LIST: hidden while a detail is open -->
      <div class="results">
        {#if searched && !loading && results.length === 0}
          <div class="empty-state">
            <span class="empty-icon">◌</span>
            <p class="empty-title">Hin argamne</p>
            <p class="empty-sub">No entry for "{query.trim()}" yet.</p>
          </div>
        {:else if results.length > 0}
          <ul>
            {#each results as word (word.id)}
              <li>
                <button class="result-btn" on:click={() => openDetail(word)}>
                  <div class="word-row">
                    <span class="term">
                      {cap(direction === 'en' ? word.english : word.afaanOromo)}
                    </span>
                    {#if word.partOfSpeech}
                      <span class="pos">{word.partOfSpeech}</span>
                    {/if}
                  </div>
                  <div class="translation">
                    {cap(direction === 'en' ? word.afaanOromo : word.english)}
                  </div>
                </button>
              </li>
            {/each}
          </ul>
        {:else if !searched}
          <div class="idle-state">
            <p>Search a word in {labels[direction].name} to begin.</p>
          </div>
        {/if}
      </div>
    {/if}
  </div>

  <footer>
    <p>Developed by <strong>Meles Zewude</strong></p>
    <p class="footer-contact">
      <!-- Replace with your real email / phone -->
      meles.zewde.tech@gmail.com &nbsp;•&nbsp; +251-92-022-7833
    </p>
  </footer>
</main>

<style>
  :global(html) {
    background: #1b2636;
  }

  main {
    max-width: 640px;
    margin: 0 auto;
    padding: 3rem 1.5rem 2rem;
    font-family: 'Nunito', system-ui, sans-serif;
    color: #edeff2;
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    box-sizing: border-box;
  }

  header {
    text-align: center;
    margin-bottom: 2rem;
  }

  h1 {
    font-size: 2rem;
    font-weight: 800;
    letter-spacing: -0.02em;
    margin: 0;
    color: #e0a458;
  }

  .subtitle {
    margin: 0.35rem 0 0;
    font-size: 0.9rem;
    color: #8fa1b8;
  }

  .translator {
    background: #22314480;
    border: 1px solid #33465c;
    border-radius: 16px;
    padding: 1.25rem;
    backdrop-filter: blur(6px);
    flex: 1;
  }

  .lang-row {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .lang-pill {
    flex: 1;
    padding: 0.6rem 1rem;
    border-radius: 999px;
    border: 1px solid #3a4d63;
    background: transparent;
    color: #9fb0c3;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .lang-pill.active {
    background: #e0a458;
    border-color: #e0a458;
    color: #1b2636;
  }

  .lang-pill:not(.active):hover {
    border-color: #e0a458;
    color: #e0a458;
  }

  .swap-btn {
    flex-shrink: 0;
    width: 38px;
    height: 38px;
    border-radius: 50%;
    border: 1px solid #3a4d63;
    background: #1b2636;
    color: #e0a458;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1), border-color 0.2s;
  }

  .swap-btn:hover {
    border-color: #e0a458;
  }

  .swap-btn.spinning {
    transform: rotate(180deg);
  }

  .search-box {
    position: relative;
    margin-bottom: 1rem;
  }

  input {
    width: 100%;
    padding: 0.9rem 2.75rem 0.9rem 1rem;
    font-size: 1.05rem;
    border-radius: 10px;
    border: 1px solid #3a4d63;
    background: #16202e;
    color: #edeff2;
    outline: none;
    box-sizing: border-box;
    transition: border-color 0.2s;
  }

  input::placeholder {
    color: #5f7488;
  }

  input:focus {
    border-color: #e0a458;
  }

  .spinner {
    position: absolute;
    right: 14px;
    top: 50%;
    width: 16px;
    height: 16px;
    margin-top: -8px;
    border: 2px solid #3a4d63;
    border-top-color: #e0a458;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* LARGE detail view */
  .detail-large {
    text-align: left;
    padding: 1rem 0.5rem 0.5rem;
    animation: fadeIn 0.2s ease;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(-6px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .back-btn {
    background: none;
    border: none;
    color: #e0a458;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    padding: 0 0 1.25rem;
  }

  .back-btn:hover {
    text-decoration: underline;
  }

  .detail-term-large {
    font-size: 2.1rem;
    font-weight: 800;
    color: #edeff2;
    display: flex;
    align-items: baseline;
    gap: 0.7rem;
    flex-wrap: wrap;
  }

  .detail-pos-large {
    font-size: 0.9rem;
    font-style: italic;
    font-weight: 400;
    color: #7c93a8;
  }

  .detail-translation-large {
    font-size: 1.6rem;
    font-weight: 700;
    color: #4fb8a6;
    margin-top: 0.5rem;
  }

  .detail-pron-large {
    font-size: 1rem;
    color: #7c93a8;
    margin-top: 0.4rem;
  }

  .detail-examples-large {
    margin-top: 1.25rem;
    padding-top: 1.25rem;
    border-top: 1px solid #2c3c4f;
  }

  .detail-examples-large p {
    margin: 0 0 0.4rem;
    font-size: 1rem;
  }

  .ex-en {
    color: #edeff2;
  }

  .ex-om {
    color: #8fa1b8;
    font-style: italic;
  }

  /* Results list */
  .results ul {
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .results li {
    border-bottom: 1px solid #2c3c4f;
  }

  .results li:last-child {
    border-bottom: none;
  }

  .result-btn {
    width: 100%;
    text-align: left;
    background: none;
    border: none;
    padding: 0.9rem 0.5rem;
    cursor: pointer;
    border-radius: 8px;
    transition: background 0.15s;
  }

  .result-btn:hover {
    background: #ffffff08;
  }

  .word-row {
    display: flex;
    align-items: baseline;
    gap: 0.5rem;
  }

  .term {
    font-size: 1.1rem;
    font-weight: 700;
    color: #edeff2;
  }

  .pos {
    font-size: 0.75rem;
    font-style: italic;
    color: #7c93a8;
  }

  .translation {
    font-size: 1rem;
    color: #4fb8a6;
    font-weight: 600;
    margin-top: 0.2rem;
  }

  .idle-state,
  .empty-state {
    text-align: center;
    padding: 2rem 1rem;
    color: #7c93a8;
  }

  .empty-icon {
    display: block;
    font-size: 1.75rem;
    color: #e0a458;
    margin-bottom: 0.5rem;
  }

  .empty-title {
    font-size: 1.1rem;
    font-weight: 700;
    color: #edeff2;
    margin: 0 0 0.25rem;
  }

  .empty-sub {
    margin: 0;
    font-size: 0.9rem;
  }

  footer {
    text-align: center;
    margin-top: 2rem;
    padding-top: 1.25rem;
    color: #5f7488;
    font-size: 0.8rem;
  }

  footer p {
    margin: 0.2rem 0;
  }

  .footer-contact {
    color: #4c5f74;
  }

  @media (max-width: 480px) {
    .lang-row {
      flex-wrap: wrap;
    }
    .lang-pill {
      flex: 1 1 40%;
    }
    .detail-term-large {
      font-size: 1.6rem;
    }
    .detail-translation-large {
      font-size: 1.3rem;
    }
  }
</style>