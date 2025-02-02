<script lang="ts">
</script>

<main class="filmfest">
    <div class="content">
        <div class="top-text">Harvard Film Festival</div>
        <div class="festival-info">
            <p class="callout">
                <strong>CALL FOR SHORT FILMS UNTIL MARCH 15TH, 2025.</strong>
            </p>
            <p>
                The Golden Arm is thrilled to announce the return of the Harvard Film Festival on April 19th,
                2025.            </p>
            <p>
                Submissions are open to undergraduates (or within one year post-graduation) from any college or university in Massachusetts. We accept original films under 30 minutes in length from any genre made since January 1st, 2023. To submit your film, please send an unlisted YouTube link to <a href="mailto:goldenarmtheater@gmail.com">goldenarmtheater@gmail.com</a> in an email with subject "[YOUR NAME]: Harvard Film Festival 2025 Submission."
            </p>
            <p>
              Direct any questions to <a href="mailto:goldenarmtheater@gmail.com">goldenarmtheater@gmail.com</a>.
            </p>

            <hr class="separator">

            <h3>Round One:</h3>
            <p>
                Films will be blindly screened by a committee of Harvard Department of Art, Film, and Visual
                Studies undergraduate and graduate students with community input at three public seances in
                March. Five finalists will be selected for the Harvard Film Festival in April. Finalists will need to
                send films in hard drive format to Harvard in late March.            </p>
            <h3>Round Two (Harvard Film Festival):</h3>
            <p>
                Finalist films will be publicly screened in the Menschel Hall Theater on Saturday April 19th. A
                panel of three Harvard film faculty members will select the winning filmmaker for a modest
                stipend and Harvard Film Festival nominal awards.            </p>
        </div>
    </div>
</main>

<style>
main.filmfest {
    padding: 1rem 2rem;
    color: #ffffff;
  }

.separator {
    margin: 2rem 0rem;
}
.top-text {
  text-align: center;
  font-size: 2rem;

  margin-bottom: 2rem;
  font-weight: bold;
  margin-top: 1rem;

}

.festival-info {
    text-align: left;
    margin-top: 3rem;
}

.callout {
    font-size: 1.2rem;
    margin-bottom: 1rem;
    color: var(--gold);
}

h3 {
    color: var(--gold);
    /* margin-top: 1rem; */
}

p {
    margin-bottom: 1rem;
    line-height: 1.6;
}

p a {
    color: var(--gold);
    text-decoration: none;
    transition: color 0.3s ease, text-decoration 0.3s ease;
}

p a:hover {
    text-decoration: underline;
}

@media screen and (max-width: 768px) {
    main.filmfest {
    margin-top: 2rem;
    padding: 2rem;
    color: #ffffff;
  }
  .top-text {
    margin-top: 3rem;

  }

  .festival-info {
    margin-top: 0rem;
  }
}
</style>
