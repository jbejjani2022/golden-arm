<script lang="ts">
    import { onMount } from 'svelte';
  
    let movie: any = null;
    let error: string = '';

    onMount(async () => {
    try {
      const response = await fetch('/api/movie/next');
      const data = await response.json();

      if (data.success) {
        movie = data.data;
      } else {
        error = 'Failed to load the next movie.';
      }
    } catch (err) {
      console.error(err);
      error = 'Something went wrong while fetching the movie data.';
    }
  });
</script>

<div style="display: flex; flex-direction: column; align-items: center; text-align: center;">
    <div class="top-text">MEET THE TEAM BEHIND THE SCREEN.</div>
    <img src="/team.jpg" alt="Golden Arm Theater team members posing for a group photo" class="team-image" />
    <div class="team-caption">
        Photography by Owen Ebose, Harvard '25.<br>
        <span class="caption-names">Left to right: Joey Bejjani, Renée Perpignan, Xander Patton, Alexandre Benoit (bottom).</span>
    </div>
    <div class="team-description" style="text-align: left; align-items: flex-start;">
        <p>
            We like movies. We like watching them, talking about them, and screening them. Sometimes we even <a href="https://xanderdraven.weebly.com/film.html" class="links">make them</a>, too.
        </p>
        <p>
            But what we like most about movies is the way they bring people together.
        </p>
        <p>
            That's why, in the fall of 2024, we revived the Golden Arm, Harvard's only student-run movie theater.
        </p>
        <p>
            Our mission is simple: to foster community through the shared joy of cinema—by screening great films, sparking conversation, and keeping the magic of the movies alive at Harvard and beyond.
        </p>
        <p>
            <br>
            Who are we? We're a team of friends in Eliot House who operate, design, engineer, bartend, and pop corn at the Golden Arm:
        </p>
        <p>
            <br>
            <a href="https://xanderdraven.weebly.com" class="links">Xander Patton</a> '26 is a student in Art, Film, and Visual Studies. Aside from co-managing the theater, he is its very own poster designer and mixologist. Find his posters and crafted drinks in our <a href="/archives" class="links">archives</a>.
        </p>
        <p>
            <a href="/about" class="links">Alexandre Benoit</a> '26 is a student in History and Literature and Art, Film, and Visual Studies. As co-manager, he plans screenings and organizes panels with special guests. He is also the Golden Arm's acclaimed Chef du Popcorn.
        </p>
        <p>
            <a href="https://github.com/jbejjani2022" class="links">Joey Bejjani</a> '26 studies Computer Science and Statistics. He engineers and maintains the <a href="https://github.com/jbejjani2022/golden-arm" class="links">software</a> behind the Golden Arm. He likes music, too—you'll find him with the <a href="https://bachsocietyorchestra.org" class="links">Bach Society Orchestra</a>.
        </p>
        <p>
            <a href="/about" class="links">Renée Perpignan</a> '26 studies Computer Science and Government. Aside from building interfaces for the users of <a href="https://goldenarmtheater.com" class="links">goldenarmtheater.com</a>, she makes your custom movie name tag when you {#if movie?.ID}<a href={`/reservations/${movie.ID}`} class="links">book a seat</a>{:else}book a seat{/if}. You'll also find her <a href="https://www.youtube.com/@reneesophia9077/videos" class="links">making music</a> and DJing.
        </p>
        <p>
            <br>
            Special thanks to our logo designer Karen Choi '26 and our merch designer Taylor Johnson '26.
        </p>
        <p>
            <br>
            Interested in joining the team? Shoot us a note at <a href="mailto:goldenarmtheater@gmail.com" class="links">goldenarmtheater@gmail.com</a>.
        </p>
        <p>
            <br>
            To many more films ahead,
            <br>
            <img src="/signature.svg" alt="Golden Arm Team signature" class="signature" />
        </p>
    </div>
</div>

<style>
    .links {
        color: var(--gold);
        text-decoration: none;
        font-weight: normal;
        transition: color 0.3s;
    }

    .top-text {
        font-family: Telegraf-Ultrabold, sans-serif;
        font-size: 30px;
        font-weight: bold;
        margin: 20px auto 40px;
        max-width: 100%;
        box-sizing: border-box;
    }

    .team-image {
        max-width: 900px;
        width: 100%;
        height: auto;
        margin-bottom: 10px;
        box-sizing: border-box;
    }

    .team-caption {
        font-size: 16px;
        margin-bottom: 40px;
        font-style: italic;
    }

    .caption-names {
        display: inline-block;
        margin-top: 6px;
    }

    .team-description {
        font-size: 20px;
        margin-bottom: 40px;
        max-width: 1050px;
        text-align: left;
        align-items: flex-start;
        box-sizing: border-box;
    }

    .signature {
        filter: invert(1);
        height: 40px;
        margin-top: 10px;
    }

    @media (max-width: 768px) {
        .top-text {
            font-size: 26px;
            margin-top: 80px;
            padding: 0 16px;
        }

        .team-description {
            padding: 0 16px;
        }
    }

    @media (max-width: 480px) {
        .top-text {
            font-size: 24px;
            margin-top: 100px;
        }
    }
</style>