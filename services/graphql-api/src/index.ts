import { createSchema, createYoga } from 'graphql-yoga';
import { createServer } from 'http';
import { config } from './config';
import { resolvers } from './resolvers';
import { typeDefs } from './schema';

const schema = createSchema({ typeDefs, resolvers });
const yoga = createYoga({ schema, graphiql: true });

const server = createServer(yoga);
server.listen(config.port, () => {
  console.log(`GraphQL API at http://localhost:${config.port}/graphql`);
});
