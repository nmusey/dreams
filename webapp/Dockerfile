FROM node:20
WORKDIR /usr/app

COPY package* .
RUN npm install
RUN npm install lightningcss --ignore-platform

COPY . .

EXPOSE 3000
CMD ["npm", "run", "dev"]
